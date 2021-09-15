package pool

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	"github.com/vesoft-inc/nebula-http-gateway/common"
	"github.com/vesoft-inc/nebula-http-gateway/service/logger"

	uuid "github.com/satori/go.uuid"
	nebula "github.com/vesoft-inc/nebula-go/v2"
	nebulaType "github.com/vesoft-inc/nebula-go/v2/nebula"
)

var (
	ConnectionClosedError = errors.New("an existing connection was forcibly closed, please check your network")
	SessionLostError      = errors.New("the connection session was lost, please connect again")
	InterruptError        = errors.New("Other statements was not executed due to this error.")
)

// Console side commands
const (
	Unknown = -1
	Param   = 1
	Params  = 2
)

type Account struct {
	username string
	password string
}

type ChannelResponse struct {
	Result *nebula.ResultSet
	Params common.ParameterMap
	Error  error
}

type ChannelRequest struct {
	Gql             string
	ResponseChannel chan ChannelResponse
	ParamList       common.ParameterList
}

type Connection struct {
	RequestChannel chan ChannelRequest
	CloseChannel   chan bool
	updateTime     int64
	parameterMap   common.ParameterMap
	account        *Account
	session        *nebula.Session
}

var connectionPool = make(map[string]*Connection)
var currentConnectionNum = 0
var connectLock sync.Mutex

func isThriftProtoError(err error) bool {
	protoErr, ok := err.(thrift.ProtocolException)
	if !ok {
		return false
	}
	if protoErr.TypeID() != thrift.UNKNOWN_PROTOCOL_EXCEPTION {
		return false
	}
	errPrefix := []string{"wsasend", "wsarecv", "write:"}
	errStr := protoErr.Error()
	for _, e := range errPrefix {
		if strings.Contains(errStr, e) {
			return true
		}
	}
	return false
}

func isThriftTransportError(err error) bool {
	if transErr, ok := err.(thrift.TransportException); ok {
		typeId := transErr.TypeID()
		if typeId == thrift.UNKNOWN_TRANSPORT_EXCEPTION || typeId == thrift.TIMED_OUT {
			if strings.Contains(transErr.Error(), "read:") {
				return true
			}
		}
	}
	return false
}

// construct Slice to nebula.NList
func Slice2Nlist(list []interface{}) (*nebulaType.NList, error) {
	sv := []*nebulaType.Value{}
	var ret nebulaType.NList
	for _, item := range list {
		nv, er := Base2Value(item)
		if er != nil {
			return nil, er
		}
		sv = append(sv, nv)
	}
	ret.Values = sv
	return &ret, nil
}

// construct map to nebula.NMap
func Map2Nmap(m map[string]interface{}) (*nebulaType.NMap, error) {
	var ret nebulaType.NMap
	kvs := map[string]*nebulaType.Value{}
	for k, v := range m {
		nv, err := Base2Value(v)
		if err != nil {
			return nil, err
		}
		kvs[k] = nv
	}
	ret.Kvs = kvs
	return &ret, nil
}

// construct go-type to nebula.Value
func Base2Value(any interface{}) (value *nebulaType.Value, err error) {
	value = nebulaType.NewValue()
	if v, ok := any.(bool); ok {
		value.BVal = &v
	} else if v, ok := any.(int); ok {
		ival := int64(v)
		value.IVal = &ival
	} else if v, ok := any.(float64); ok {
		if v == float64(int64(v)) {
			iv := int64(v)
			value.IVal = &iv
		} else {
			value.FVal = &v
		}
	} else if v, ok := any.(float32); ok {
		if v == float32(int64(v)) {
			iv := int64(v)
			value.IVal = &iv
		} else {
			fval := float64(v)
			value.FVal = &fval
		}
	} else if v, ok := any.(string); ok {
		value.SVal = []byte(v)
	} else if any == nil {
		nval := nebulaType.NullType___NULL__
		value.NVal = &nval
	} else if v, ok := any.([]interface{}); ok {
		nv, er := Slice2Nlist([]interface{}(v))
		if er != nil {
			err = er
		}
		value.LVal = nv
	} else if v, ok := any.(map[string]interface{}); ok {
		nv, er := Map2Nmap(map[string]interface{}(v))
		if er != nil {
			err = er
		}
		value.MVal = nv
	} else {
		// unsupport other Value type, use this function carefully
		err = fmt.Errorf("Only support convert boolean/float/int/string/map/list to nebula.Value but %T", any)
	}
	return
}

func isCmd(query string) (isLocal bool, localCmd int, args []string) {
	isLocal = false
	localCmd = Unknown
	plain := strings.TrimSpace(query)
	if len(plain) < 1 || plain[0] != ':' {
		return
	}
	isLocal = true
	words := strings.Fields(plain[1:])
	localCmdName := words[0]
	switch strings.ToLower(localCmdName) {
	case "param":
		{
			localCmd = Param
			args = []string{plain}
		}
	case "params":
		{
			localCmd = Params
			args = []string{plain}
		}
	}
	return
}

func executeCmd(parameterList common.ParameterList, parameterMap *common.ParameterMap) (showMap common.ParameterMap, err error) {
	for _, v := range parameterList {
		// convert interface{} to nebula.Value
		if isLocal, cmd, args := isCmd(v); isLocal {
			switch cmd {
			case Param:
				if len(args) == 1 {
					err = defineParams(args[0], parameterMap)
				}
				if err != nil {
					return nil, err
				}
			case Params:
				if len(args) == 1 {
					showMap, err = ListParams(args[0], parameterMap)
				}
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return showMap, nil
}

func defineParams(args string, parameterMap *common.ParameterMap) (err error) {
	argsRewritten := strings.Replace(args, "'", "\"", -1)
	reg := regexp.MustCompile(`^\s*:param\s+(\S+)\s*=>(.*)$`)
	if reg == nil {
		err = errors.New("invalid regular expression")
		return
	}
	matchResult := reg.FindAllStringSubmatch(argsRewritten, -1)
	if len(matchResult) != 1 || len(matchResult[0]) != 3 {
		err = errors.New("Set params failed. Wrong local command format (" + reg.String() + ") ")
		return
	}
	/*
	 * :param p1=> -> [":param p1=>",":p1",""]
	 * :param p2=>3 -> [":param p2=>3",":p2","3"]
	 */
	paramKey := matchResult[0][1]
	paramValue := matchResult[0][2]
	if len(paramValue) == 0 {
		delete((*parameterMap), paramKey)
	} else {
		paramsWithGoType := make(common.ParameterMap)
		param := "{\"" + paramKey + "\"" + ":" + paramValue + "}"
		err = json.Unmarshal([]byte(param), &paramsWithGoType)
		if err != nil {
			return
		}
		for k, v := range paramsWithGoType {
			(*parameterMap)[k] = v
		}
	}
	return nil
}

func ListParams(args string, parameterMap *common.ParameterMap) (showMap common.ParameterMap, err error) {
	reg := regexp.MustCompile(`^\s*:params\s*(\S*)\s*$`)
	paramsWithGoType := make(common.ParameterMap)
	if reg == nil {
		err = errors.New("invalid regular expression")
		return
	}
	matchResult := reg.FindAllStringSubmatch(args, -1)
	if len(matchResult) != 1 {
		err = errors.New("Set params failed. Wrong local command format " + reg.String() + ") ")
		return
	}
	res := matchResult[0]
	/*
	 * :params -> [":params",""]
	 * :params p1 -> ["params","p1"]
	 */
	if len(res) != 2 {
		return
	} else {
		paramKey := matchResult[0][1]
		if len(paramKey) == 0 {
			for k, v := range *parameterMap {
				paramsWithGoType[k] = v
			}
		} else {
			if paramValue, ok := (*parameterMap)[paramKey]; ok {
				paramsWithGoType[paramKey] = paramValue
			} else {
				err = errors.New("Unknown parameter: " + paramKey)
			}
		}
	}
	return paramsWithGoType, nil
}

func NewConnection(address string, port int, username string, password string) (nsid string, err error) {
	connectLock.Lock()
	defer connectLock.Unlock()
	// Initialize logger
	var httpgatewayLog = logger.HttpGatewayLogger{}
	hostAddress := nebula.HostAddress{Host: address, Port: port}
	hostList := []nebula.HostAddress{hostAddress}
	poolConfig := nebula.GetDefaultConf()
	// Initialize connectin pool
	pool, err := nebula.NewConnectionPool(hostList, poolConfig, httpgatewayLog)
	if err != nil {
		return "", err
	}
	err = pool.Ping(hostList[0], 5000*time.Millisecond)
	if err != nil {
		return "", err
	}

	// Create session
	session, err := pool.GetSession(username, password)
	if err == nil {
		nsid = uuid.NewV4().String()
		connectionPool[nsid] = &Connection{
			RequestChannel: make(chan ChannelRequest),
			CloseChannel:   make(chan bool),
			updateTime:     time.Now().Unix(),
			session:        session,
			parameterMap:   make(common.ParameterMap),
			account: &Account{
				username: username,
				password: password,
			},
		}
		currentConnectionNum++

		// Make a goroutine to deal with concurrent requests from each connection
		go func() {
			connection := connectionPool[nsid]
			for {
				select {
				case request := <-connection.RequestChannel:
					func() {
						defer func() {
							if err := recover(); err != nil {
								common.LogPanic(err)
								request.ResponseChannel <- ChannelResponse{
									Result: nil,
									Error:  SessionLostError,
								}
							}
						}()
						showMap := make(common.ParameterMap)
						if len(request.ParamList) > 0 {
							showMap, err = executeCmd(request.ParamList, &connection.parameterMap)
							if err != nil {
								if len(request.Gql) > 0 {
									err = errors.New(err.Error() + InterruptError.Error())
								}
								request.ResponseChannel <- ChannelResponse{
									Result: nil,
									Params: showMap,
									Error:  err,
								}
								return
							}
						}

						if len(request.Gql) > 0 {
							params := make(map[string]*nebulaType.Value)
							for k, v := range connection.parameterMap {
								value, paramError := Base2Value(v)
								if paramError != nil {
									err = paramError
								}
								params[k] = value
							}
							response, err := connection.session.ExecuteWithParameter(request.Gql, params)
							if err != nil && (isThriftProtoError(err) || isThriftTransportError(err)) {
								err = ConnectionClosedError
							}
							request.ResponseChannel <- ChannelResponse{
								Result: response,
								Params: showMap,
								Error:  err,
							}
						} else {
							request.ResponseChannel <- ChannelResponse{
								Result: nil,
								Params: showMap,
								Error:  nil,
							}
						}
					}()
				case <-connection.CloseChannel:
					connection.session.Release()
					connectLock.Lock()
					delete(connectionPool, nsid)
					currentConnectionNum--
					connectLock.Unlock()
					// Exit loop
					return
				}
			}
		}()
		return nsid, err
	}
	return "", err
}

func Disconnect(nsid string) {
	connection := connectionPool[nsid]
	if connection != nil {
		connection.session.Release()
		delete(connectionPool, nsid)
	}
}

func GetConnection(nsid string) (connection *Connection, err error) {
	connectLock.Lock()
	defer connectLock.Unlock()

	if connection, ok := connectionPool[nsid]; ok {
		connection.updateTime = time.Now().Unix()
		return connection, nil
	}
	return nil, errors.New("connection refused for being released")
}
