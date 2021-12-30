package v2_6

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_6/graph"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

var (
	_ types.GraphClientDriver = (*defaultGraphClient)(nil)
)

type (
	defaultGraphClient struct {
		graph *graph.GraphServiceClient
	}
)

func newGraphClient(transport thrift.Transport, pf thrift.ProtocolFactory) types.GraphClientDriver {
	return &defaultGraphClient{
		graph: graph.NewGraphServiceClientFactory(transport, pf),
	}
}

func (c *defaultGraphClient) Open() error {
	return c.graph.Open()
}

func (c *defaultGraphClient) VerifyClientVersion() error {
	req := graph.NewVerifyClientVersionReq()
	resp, err := c.graph.VerifyClientVersion(req)
	if err != nil {
		return err
	}
	return codeErrorIfHappened(resp.ErrorCode, resp.ErrorMsg)
}

func (c *defaultGraphClient) Authenticate(username, password string) (types.AuthResponse, error) {
	resp, err := c.graph.Authenticate([]byte(username), []byte(password))
	if err != nil {
		return nil, err
	}

	if err = codeErrorIfHappened(resp.ErrorCode, resp.ErrorMsg); err != nil {
		return nil, err
	}
	return newAuthResponseWrapper(resp), nil
}

func (c *defaultGraphClient) Signout(sessionId int64) (err error) {
	return c.graph.Signout(sessionId)
}

func (c *defaultGraphClient) Execute(sessionId int64, stmt []byte) (types.ExecutionResponse, error) {
	resp, err := c.graph.Execute(sessionId, stmt)
	if err != nil {
		return nil, err
	}

	if err = codeErrorIfHappened(resp.ErrorCode, resp.ErrorMsg); err != nil {
		return nil, err
	}
	return newEexecutionResponseWrapper(resp), nil
}

func (c *defaultGraphClient) ExecuteJson(sessionId int64, stmt []byte) ([]byte, error) {
	return c.graph.ExecuteJson(sessionId, stmt)
}

func (c *defaultGraphClient) Close() error {
	if c.graph != nil {
		if err := c.graph.Close(); err != nil {
			return err
		}
	}
	return nil
}
