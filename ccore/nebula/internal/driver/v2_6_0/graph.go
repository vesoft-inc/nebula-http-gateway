package v2_6_0

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	nerrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
	graph2_6_0 "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_6_0/graph"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

var (
	_ types.GraphClientDriver = (*defaultGraphClient)(nil)
)

type (
	defaultGraphClient struct {
		graph *graph2_6_0.GraphServiceClient
	}

	authResponse struct {
		*graph2_6_0.AuthResponse
	}

	executionResponse struct {
		*graph2_6_0.ExecutionResponse
	}
)

func newGraphClient(transport thrift.Transport, pf thrift.ProtocolFactory) types.GraphClientDriver {
	return &defaultGraphClient{
		graph: graph2_6_0.NewGraphServiceClientFactory(transport, pf),
	}
}

func (c *defaultGraphClient) Open() error {
	return c.graph.Open()
}

func (c *defaultGraphClient) Authenticate(username, password string) (types.AuthResponse, error) {
	resp, err := c.graph.Authenticate([]byte(username), []byte(password))
	if err != nil {
		return nil, err
	}
	return &authResponse{AuthResponse: resp}, nil
}

func (c *defaultGraphClient) Signout(sessionId int64) (err error) {
	return c.graph.Signout(sessionId)
}

func (c *defaultGraphClient) Execute(sessionId int64, stmt []byte) (types.ExecutionResponse, error) {
	resp, err := c.graph.Execute(sessionId, stmt)
	if err != nil {
		return nil, err
	}

	return &executionResponse{ExecutionResponse: resp}, nil
}

func (c *defaultGraphClient) ExecuteJson(sessionId int64, stmt []byte) ([]byte, error) {
	return c.graph.ExecuteJson(sessionId, stmt)
}

func (c *defaultGraphClient) Close() error{
	if c.graph != nil {
		if err := c.graph.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (r *authResponse) ErrorCode() nerrors.ErrorCode {
	return nerrors.ErrorCode(r.AuthResponse.ErrorCode)
}

func (r *authResponse) ErrorMsg() string {
	return string(r.AuthResponse.ErrorMsg)
}

func (r *authResponse) SessionID() *int64 {
	return r.AuthResponse.SessionID
}

func (r *executionResponse) TODO() {
	// TODO: the others
}
