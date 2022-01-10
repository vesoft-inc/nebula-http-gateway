package v2_5

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	nerrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_5/graph"
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

func (c *defaultGraphClient) VerifyClientVersion() error {
	// v2.5 is not support verify client version, and it's the lowest version, so return not error.
	return nil
}

func (c *defaultGraphClient) Open() error {
	return c.graph.Open()
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
	return newExecutionResponseWrapper(resp), nil
}

func (c *defaultGraphClient) ExecuteJson(sessionId int64, stmt []byte) ([]byte, error) {
	return c.graph.ExecuteJson(sessionId, stmt)
}

func (c *defaultGraphClient) ExecuteWithParameter(sessionId int64, stmt []byte, params map[string]types.Value) (types.ExecutionResponse, error) {
	if params == nil {
		return c.Execute(sessionId, stmt)
	}
	return nil, nerrors.ErrUnsupported
}

func (c *defaultGraphClient) Close() error {
	if c.graph != nil {
		if err := c.graph.Close(); err != nil {
			return err
		}
	}
	return nil
}
