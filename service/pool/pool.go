package pool

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	uuid "github.com/satori/go.uuid"
	nebula "github.com/vesoft-inc/nebula-go"
	common "nebula-http-gateway/utils"
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

const (
	maxConnectionNum  = 200
	secondsOfHalfHour = int64(30 * 60)
)

var connectionPool = make(map[string]*Connection)
var currentConnectionNum = 0
var connectLock sync.Mutex

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
		return "", errors.New("Fail to initialize the connection pool")
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
					func(gql string) {
						defer func() {
							if err := recover(); err != nil {
								common.LogPanic(err)
								request.ResponseChannel <- ChannelResponse{
									Result: nil,
									Error:  SessionLostError,
								}
							}
						}()
						response, err := connection.session.Execute(gql)
						if protoErr, ok := err.(thrift.ProtocolException); ok && protoErr != nil &&
							protoErr.TypeID() == thrift.UNKNOWN_PROTOCOL_EXCEPTION {
							if strings.Contains(protoErr.Error(), "wsasend") ||
								strings.Contains(protoErr.Error(), "wsarecv") ||
								strings.Contains(protoErr.Error(), "write:") {
								err = ConnectionClosedError
							}
						}
						if transErr, ok := err.(thrift.TransportException); ok && transErr != nil {
							if transErr.TypeID() == thrift.UNKNOWN_TRANSPORT_EXCEPTION ||
								transErr.TypeID() == thrift.TIMED_OUT {
								if strings.Contains(transErr.Error(), "read:") {
									err = ConnectionClosedError
								}
							}
						}
						request.ResponseChannel <- ChannelResponse{
							Result: response,
							Error:  err,
						}
					}(request.Gql)
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

	connection, ok := connectionPool[nsid]
	if ok {
		connection.updateTime = time.Now().Unix()
		return connection, nil
	}
	return nil, errors.New("connection refused for being released")
}

func recoverConnections() {
	nowTimeStamps := time.Now().Unix()
	for _, connection := range connectionPool {
		// release connection if not use over 30minutes
		if nowTimeStamps-connection.updateTime > secondsOfHalfHour {
			connection.CloseChannel <- true
		}
	}
}
