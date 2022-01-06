package pool

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	"github.com/vesoft-inc/nebula-http-gateway/common"
	"github.com/vesoft-inc/nebula-http-gateway/service/logger"

	uuid "github.com/satori/go.uuid"
	nebula "github.com/vesoft-inc/nebula-go/v2"
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
		localCmd = Param
		args = []string{plain}
	case "params":
		localCmd = Params
		args = []string{plain}
	}
	return
}

func executeCmd(parameterList common.ParameterList, parameterMap *common.ParameterMap) (showMap common.ParameterMap, err error) {
	tempMap := make(common.ParameterMap)
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
					err = ListParams(args[0], &tempMap, parameterMap)
				}
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return tempMap, nil
}

func defineParams(args string, parameterMap *common.ParameterMap) (err error) {
	argsRewritten := strings.Replace(args, "'", "\"", -1)
	reg := regexp.MustCompile(`(?i)^\s*:param\s+(\S+)\s*=>(.*)$`)
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

func ListParams(args string, tmpParameter *common.ParameterMap, sessionMap *common.ParameterMap) (err error) {
	reg := regexp.MustCompile(`(?i)^\s*:params\s*(\S*)\s*$`)
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
	}
	paramKey := matchResult[0][1]
	if len(paramKey) == 0 {
		for k, v := range *sessionMap {
			(*tmpParameter)[k] = v
		}
	} else {
		if paramValue, ok := (*sessionMap)[paramKey]; ok {
			(*tmpParameter)[paramKey] = paramValue
		} else {
			err = errors.New("Unknown parameter: " + paramKey)
		}
	}
	return nil
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
							response, err := connection.session.ExecuteWithParameter(request.Gql, connection.parameterMap)
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
