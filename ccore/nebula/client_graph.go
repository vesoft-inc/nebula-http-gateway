package nebula

import (
	"fmt"
	nerrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
)

type (
	GraphClient interface {
		Open() error
		Execute(stmt []byte) (ExecutionResponse, error)
		ExecuteJson(stmt []byte) ([]byte, error)
		Close() error
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
	err := c.graph.open(c.driver)
	if err != nil {
		return err
	}

	resp, err := c.graph.VerifyClientVersion()
	if err != nil {
		return fmt.Errorf("failed to verify client version: %s", err.Error())
	}
	if resp != nil && resp.ErrorCode != nerrors.ErrorCode_SUCCEEDED {
		return fmt.Errorf("incompatible version between client and server: %s", string(resp.ErrorMsg))
	}
	return nil
}

func (c *defaultGraphClient) Execute(stmt []byte) (ExecutionResponse, error) {
	return c.graph.Execute(c.graph.sessionId, stmt)
}

func (c *defaultGraphClient) ExecuteJson(stmt []byte) ([]byte, error) {
	if err := c.graph.open(c.driver); err != nil {
		return nil, err
	}

	return c.graph.ExecuteJson(c.graph.sessionId, stmt)
}

func (c *defaultGraphClient) Close() error {
	return c.graph.close()
}
