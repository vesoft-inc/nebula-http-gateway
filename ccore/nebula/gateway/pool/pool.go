package pool

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	uuid "github.com/satori/go.uuid"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/wrapper"
)

var (
	ConnectionClosedError = errors.New("an existing connection was forcibly closed, please check your network")
	SessionLostError      = errors.New("the connection session was lost, please connect again")
	InterruptError        = errors.New("other statements was not executed due to this error")
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
	Result *wrapper.ResultSet
	Params types.ParameterMap
	Msg    interface{}
	Error  error
}

type ChannelRequest struct {
	Gql             string
	ResponseChannel chan ChannelResponse
	ParamList       types.ParameterList
}

type Client struct {
	graphClient    nebula.GraphClient
	RequestChannel chan ChannelRequest
	CloseChannel   chan bool
	updateTime     int64
	parameterMap   types.ParameterMap
	account        *Account
}

var (
	clientPool       = make(map[string]*Client)
	currentClientNum = 0
	clientMux        sync.Mutex

	ClientNotExistedError = errors.New("get client error: client not existed, session expired")
)

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
func Slice2Nlist(factory types.FactoryDriver, list []interface{}) (types.NList, error) {
	sv := make([]types.Value, 0, len(list))
	nListBuilder := factory.NewNListBuilder()
	for _, item := range list {
		nv, er := Base2Value(factory, item)
		if er != nil {
			return nil, er
		}
		sv = append(sv, nv)
	}
	nListBuilder.Values(sv)
	return nListBuilder.Build(), nil
}

// construct map to nebula.NMap
func Map2Nmap(factory types.FactoryDriver, m map[string]interface{}) (types.NMap, error) {
	nMapBuilder := factory.NewNMapBuilder()
	kvs := map[string]types.Value{}
	for k, v := range m {
		nv, err := Base2Value(factory, v)
		if err != nil {
			return nil, err
		}
		kvs[k] = nv
	}
	nMapBuilder.Kvs(kvs)
	return nMapBuilder.Build(), nil
}

// construct go-type to nebula.Value
func Base2Value(factory types.FactoryDriver, any interface{}) (types.Value, error) {
	var err error
	valueBuilder := factory.NewValueBuilder()
	if v, ok := any.(bool); ok {
		valueBuilder.BVal(&v)
	} else if v, ok := any.(int); ok {
		ival := int64(v)
		valueBuilder.IVal(&ival)
	} else if v, ok := any.(float64); ok {
		if v == float64(int64(v)) {
			iv := int64(v)
			valueBuilder.IVal(&iv)
		} else {
			valueBuilder.FVal(&v)
		}
	} else if v, ok := any.(float32); ok {
		if v == float32(int64(v)) {
			iv := int64(v)
			valueBuilder.IVal(&iv)
		} else {
			fval := float64(v)
			valueBuilder.FVal(&fval)
		}
	} else if v, ok := any.(string); ok {
		valueBuilder.SVal([]byte(v))
	} else if any == nil {
		nval := types.NullType___NULL__
		valueBuilder.NVal(&nval)
	} else if v, ok := any.([]interface{}); ok {
		nv, er := Slice2Nlist(factory, []interface{}(v))
		if er != nil {
			err = er
		}
		valueBuilder.LVal(nv)
	} else if v, ok := any.(map[string]interface{}); ok {
		nv, er := Map2Nmap(factory, map[string]interface{}(v))
		if er != nil {
			err = er
		}
		valueBuilder.MVal(nv)
	} else {
		// unsupport other Value type, use this function carefully
		err = fmt.Errorf("Only support convert boolean/float/int/string/map/list to nebula.Value but %T", any)
	}
	return valueBuilder.Build(), err
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

func executeCmd(parameterList types.ParameterList, parameterMap types.ParameterMap) (showMap types.ParameterMap, err error) {
	tempMap := make(types.ParameterMap)
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
					err = ListParams(args[0], tempMap, parameterMap)
				}
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return tempMap, nil
}

func defineParams(args string, parameterMap types.ParameterMap) (err error) {
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
		delete(parameterMap, paramKey)
	} else {
		paramsWithGoType := make(types.ParameterMap)
		param := "{\"" + paramKey + "\"" + ":" + paramValue + "}"
		err = json.Unmarshal([]byte(param), &paramsWithGoType)
		if err != nil {
			return
		}
		for k, v := range paramsWithGoType {
			parameterMap[k] = v
		}
	}
	return nil
}

func ListParams(args string, tmpParameter types.ParameterMap, sessionMap types.ParameterMap) (err error) {
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
		for k, v := range sessionMap {
			tmpParameter[k] = v
		}
	} else {
		if paramValue, ok := sessionMap[paramKey]; ok {
			tmpParameter[paramKey] = paramValue
		} else {
			err = errors.New("Unknown parameter: " + paramKey)
		}
	}
	return nil
}

func NewClient(address string, port int, username string, password string, version nebula.Version) (ncid string, err error) {
	clientMux.Lock()
	defer clientMux.Unlock()

	host := strings.Join([]string{address, strconv.Itoa(port)}, ":")
	c, err := nebula.NewGraphClient([]string{host}, username, password, nebula.WithVersion(version))
	if err != nil {
		return "", err
	}
	if err := c.Open(); err != nil {
		return "", err
	}

	ncid = uuid.NewV4().String()

	client := &Client{
		graphClient:    c,
		RequestChannel: make(chan ChannelRequest),
		CloseChannel:   make(chan bool),
		updateTime:     time.Now().Unix(),
		parameterMap:   make(types.ParameterMap),
		account: &Account{
			username: username,
			password: password,
		},
	}
	clientPool[ncid] = client
	currentClientNum++

	// Make a goroutine to deal with concurrent requests from each connection
	go handleRequest(ncid)

	return ncid, err
}

func handleRequest(ncid string) {
	var err error
	client := clientPool[ncid]
	for {
		select {
		case request := <-client.RequestChannel:
			func() {
				defer func() {
					if err := recover(); err != nil {
						request.ResponseChannel <- ChannelResponse{
							Result: nil,
							Msg:    err,
							Error:  SessionLostError,
						}
					}
				}()
				showMap := make(types.ParameterMap)
				if request.ParamList != nil && len(request.ParamList) > 0 {
					showMap, err = executeCmd(request.ParamList, client.parameterMap)
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
					params := make(map[string]types.Value)
					for k, v := range client.parameterMap {
						value, paramError := Base2Value(client.graphClient.Factory(), v)
						if paramError != nil {
							err = paramError
						}
						params[k] = value
					}
					authResp, err := client.graphClient.Authenticate(client.account.username, client.account.password)
					if err != nil {
						err = errors.New(err.Error() + InterruptError.Error())
						request.ResponseChannel <- ChannelResponse{
							Result: nil,
							Error:  err,
						}
						return
					}

					execResponse, err := client.graphClient.ExecuteWithParameter([]byte(request.Gql), params)
					if err != nil && (isThriftProtoError(err) || isThriftTransportError(err)) {
						err = ConnectionClosedError
						request.ResponseChannel <- ChannelResponse{
							Result: nil,
							Error:  err,
						}
						return
					}

					res, err := wrapper.GenResultSet(execResponse, client.graphClient.Factory(), authResp.GetTimezoneInfo())
					if err != nil {
						err = errors.New(err.Error() + InterruptError.Error())
					}
					request.ResponseChannel <- ChannelResponse{
						Result: res,
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
		case <-client.CloseChannel:
			clientMux.Lock()
			_ = client.graphClient.Close()
			currentClientNum--
			delete(clientPool, ncid)
			clientMux.Unlock()
			return // Exit loop
		}
	}
}

func Close(ncid string) {
	clientMux.Lock()
	defer clientMux.Unlock()

	if client, ok := clientPool[ncid]; ok {
		_ = client.graphClient.Close()
		currentClientNum--
		delete(clientPool, ncid)
	}
}

func GetClient(ncid string) (*Client, error) {
	clientMux.Lock()
	defer clientMux.Unlock()

	if client, ok := clientPool[ncid]; ok {
		return client, nil
	}

	return nil, ClientNotExistedError
}
