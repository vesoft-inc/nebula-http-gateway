package nebula

import (
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

type (
	GraphClient interface {
		Open() error
		Authenticate(username, password string) (AuthResponse, error)
		Execute(stmt []byte) (ExecutionResponse, error)
		ExecuteJson(stmt []byte) ([]byte, error)
		Close() error
		Factory() Factory
		Version() Version
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
