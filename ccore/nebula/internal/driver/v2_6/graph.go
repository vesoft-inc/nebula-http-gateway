package v2_6

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	nerrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
	graph2_6 "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_6/graph"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

var (
	_ types.GraphClientDriver = (*defaultGraphClient)(nil)
)

type (
	defaultGraphClient struct {
		graph *graph2_6.GraphServiceClient
	}

	authResponse struct {
		*graph2_6.AuthResponse
	}

	executionResponse struct {
		*graph2_6.ExecutionResponse
	}
)

func newGraphClient(transport thrift.Transport, pf thrift.ProtocolFactory) types.GraphClientDriver {
	return &defaultGraphClient{
		graph: graph2_6.NewGraphServiceClientFactory(transport, pf),
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

func (c *defaultGraphClient) Close() error {
	if c.graph != nil {
		if err := c.graph.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (c *defaultGraphClient) VerifyClientVersion() (*types.VerifyClientVersionResp, error) {
	req := graph2_6.NewVerifyClientVersionReq()
	resp, err := c.graph.VerifyClientVersion(req)
	if err != nil {
		return &types.VerifyClientVersionResp{
			ErrorCode: nerrors.ErrorCode(resp.GetErrorCode()),
			ErrorMsg:  resp.GetErrorMsg(),
		}, err
	}
	return &types.VerifyClientVersionResp{
		ErrorCode: nerrors.ErrorCode_SUCCEEDED,
	}, nil
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

func (r *executionResponse) GetErrorCode() nerrors.ErrorCode {
	return nerrors.ErrorCode(r.ExecutionResponse.GetErrorCode())
}

func (r *executionResponse) GetLatencyInUs() int64 {
	return int64(r.ExecutionResponse.GetLatencyInUs())
}

func (r *executionResponse) GetData() types.DataSet {
	dataset := dataSetWrapper(r.ExecutionResponse.GetData())
	return dataset
}

func (r *executionResponse) GetSpaceName() []byte {
	return r.ExecutionResponse.GetSpaceName()
}

func (r *executionResponse) GetErrorMsg() []byte {
	return r.ExecutionResponse.GetErrorMsg()
}

func (r *executionResponse) GetPlanDesc() types.PlanDescription {
	planDesc := planDescriptionWrapper(r.PlanDesc)
	return planDesc
}

func (r *executionResponse) GetComment() []byte {
	return r.ExecutionResponse.GetComment()
}

func (r *executionResponse) IsSetData() bool {
	return r.ExecutionResponse.IsSetData()
}

func (r *executionResponse) IsSetSpaceName() bool {
	return r.ExecutionResponse.IsSetSpaceName()
}

func (r *executionResponse) IsSetErrorMsg() bool {
	return r.ExecutionResponse.IsSetErrorMsg()
}

func (r *executionResponse) IsSetPlanDesc() bool {
	return r.ExecutionResponse.IsSetPlanDesc()
}

func (r *executionResponse) IsSetComment() bool {
	return r.ExecutionResponse.IsSetComment()
}

func (r *executionResponse) String() string {
	return r.ExecutionResponse.String()
}
