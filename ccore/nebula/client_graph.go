package nebula

import (
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/wrapper"
)

type (
	GraphClient interface {
		Open() error
		Authenticate(username, password string) (AuthResponse, error)
		Execute(stmt []byte) (ExecutionResponse, error)
		ExecuteJson(stmt []byte) ([]byte, error)
		ExecuteWithParameter(stmt []byte, params types.ParameterMap) (ExecutionResponse, error)
		Close() error
		Factory() Factory
		Version() Version
		GetTimezoneInfo() types.TimezoneInfo
	}

	defaultGraphClient defaultClient
)

func NewGraphClient(endpoints []string, username, password string, opts ...Option) (GraphClient, error) {
	c, err := NewClient(ConnectionInfo{
		GraphEndpoints: endpoints,
		GraphAccount: Account{
			Username: username,
			Password: password,
		},
	}, opts...)
	if err != nil {
		return nil, err
	}
	return c.Graph(), nil
}

func (c *defaultGraphClient) GetTimezoneInfo() types.TimezoneInfo {
	return c.graph.GetTimezoneInfo()
}

func (c *defaultGraphClient) Open() error {
	return c.defaultClient().initDriver(func(driver types.Driver) error {
		return c.graph.open(driver)
	})
}

func (c *defaultGraphClient) Authenticate(username, password string) (AuthResponse, error) {
	return c.graph.Authenticate(username, password)
}

func (c *defaultGraphClient) Execute(stmt []byte) (ExecutionResponse, error) {
	return c.graph.Execute(c.graph.sessionId, stmt)
}

func (c *defaultGraphClient) ExecuteJson(stmt []byte) ([]byte, error) {
	return c.graph.ExecuteJson(c.graph.sessionId, stmt)
}

func (c *defaultGraphClient) ExecuteWithParameter(stmt []byte, params types.ParameterMap) (ExecutionResponse, error) {
	if len(params) == 0 {
		return c.Execute(stmt)
	}
	// wrap the map of interface{} to map of types.Value
	paramsMap := make(map[string]types.Value)
	for k, v := range params {
		nv, er := wrapper.WrapValue(v, c.Factory())
		if er != nil {
			return nil, er
		}
		paramsMap[k] = nv
	}
	return c.graph.ExecuteWithParameter(c.graph.sessionId, stmt, paramsMap)
}

func (c *defaultGraphClient) Close() error {
	return c.graph.close()
}

func (c *defaultGraphClient) Factory() Factory {
	return c.defaultClient().Factory()
}

func (c *defaultGraphClient) Version() Version {
	return c.defaultClient().Version()
}

func (c *defaultGraphClient) defaultClient() *defaultClient {
	return (*defaultClient)(c)
}
