package pool

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/wrapper"
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
	Result *wrapper.ResultSet
	Msg    interface{}
	Error  error
}

type ChannelRequest struct {
	Gql             string
	ResponseChannel chan ChannelResponse
}

type Client struct {
	graphClient    nebula.GraphClient
	RequestChannel chan ChannelRequest
	CloseChannel   chan bool
	updateTime     int64
	account        *Account
}

var (
	clientPool       = make(map[string]*Client)
	currentClientNum = 0
	clientMux        sync.Mutex

	ClientNotExistedError = errors.New("get client error: client not existed")
)

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
				response, err := client.graphClient.Execute([]byte(request.Gql))
				if err != nil {
					err = ConnectionClosedError
				}
				res, err := wrapper.GenResultSet(response)
				if err != nil {
					err = ConnectionClosedError
				}
				request.ResponseChannel <- ChannelResponse{
					Result: res,
					Error:  err,
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

func (c *Client) Execute(gql string) (nebula.ExecutionResponse, error) {
	return c.graphClient.Execute([]byte(gql))
}
