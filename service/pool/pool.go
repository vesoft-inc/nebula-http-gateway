package pool

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	"github.com/vesoft-inc/nebula-http-gateway/common"

	uuid "github.com/satori/go.uuid"
	nebula "github.com/vesoft-inc/nebula-go/v2"
)

var (
	ConnectionClosedError = errors.New("an existing connection was forcibly closed, please check your network")
	SessionLostError      = errors.New("the connection session was lost, please connect again")
)

type Account struct {
	username string
	password string
}

type ChannelResponse struct {
	Result *nebula.ResultSet
	Error  error
}

type ChannelRequest struct {
	Gql             string
	ResponseChannel chan ChannelResponse
}

type Connection struct {
	RequestChannel chan ChannelRequest
	CloseChannel   chan bool
	updateTime     int64
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

func NewConnection(address string, port int, username string, password string) (nsid string, err error) {
	connectLock.Lock()
	defer connectLock.Unlock()
	// Initialize logger
	var nebulaLog = nebula.DefaultLogger{}
	hostAddress := nebula.HostAddress{Host: address, Port: port}
	hostList := []nebula.HostAddress{hostAddress}
	poolConfig := nebula.GetDefaultConf()
	// Initialize connectin pool
	pool, err := nebula.NewConnectionPool(hostList, poolConfig, nebulaLog)
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
					defer func() {
						if err := recover(); err != nil {
							common.LogPanic(err)
							request.ResponseChannel <- ChannelResponse{
								Result: nil,
								Error:  SessionLostError,
							}
						}
					}()
					response, err := connection.session.Execute(request.Gql)
					if err != nil && (isThriftProtoError(err) || isThriftTransportError(err)) {
						err = ConnectionClosedError
					}
					request.ResponseChannel <- ChannelResponse{
						Result: response,
						Error:  err,
					}
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
