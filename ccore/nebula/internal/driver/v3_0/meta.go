package v3_0

import (
	"net"
	"strconv"

	"github.com/facebook/fbthrift/thrift/lib/go/thrift"

	nerrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
	nthrift "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v3_0"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v3_0/meta"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

var (
	_ types.MetaClientDriver = (*defaultMetaClient)(nil)
)

type (
	defaultMetaClient struct {
		meta *meta.MetaServiceClient
	}
)

func newMetaClient(transport thrift.Transport, pf thrift.ProtocolFactory) types.MetaClientDriver {
	return &defaultMetaClient{
		meta: meta.NewMetaServiceClientFactory(transport, pf),
	}
}

func (c *defaultMetaClient) Open() error {
	return c.meta.Open()
}

func (c *defaultMetaClient) VerifyClientVersion() error {
	req := meta.NewVerifyClientVersionReq()
	resp, err := c.meta.VerifyClientVersion(req)
	if err != nil {
		return err
	}
	return codeErrorIfHappened(resp.Code, resp.ErrorMsg)
}

func (c *defaultMetaClient) Close() error {
	if c.meta != nil {
		if err := c.meta.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (c *defaultMetaClient) AddHosts(endpoints []string) error {
	hostsToAdd := make([]*nthrift.HostAddr, 0, len(endpoints))
	for _, ep := range endpoints {
		host, portStr, err := net.SplitHostPort(ep)
		if err != nil {
			return err
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}

		hostsToAdd = append(hostsToAdd, &nthrift.HostAddr{
			Host: host,
			Port: nthrift.Port(port),
		})
	}

	req := &meta.AddHostsReq{
		Hosts: hostsToAdd,
	}
	resp, err := c.meta.AddHosts(req)
	if err != nil {
		return err
	}
	return codeErrorIfHappened(resp.Code, nil)
}

func (c *defaultMetaClient) DropHosts(endpoints []string) error {
	hostsToDrop := make([]*nthrift.HostAddr, 0, len(endpoints))
	for _, ep := range endpoints {
		host, portStr, err := net.SplitHostPort(ep)
		if err != nil {
			return err
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			return err
		}

		hostsToDrop = append(hostsToDrop, &nthrift.HostAddr{
			Host: host,
			Port: nthrift.Port(port),
		})
	}

	req := &meta.DropHostsReq{
		Hosts: hostsToDrop,
	}
	resp, err := c.meta.DropHosts(req)
	if err != nil {
		return err
	}
	return codeErrorIfHappened(resp.Code, nil)
}

func (c *defaultMetaClient) ListSpaces() (types.Spaces, error) {
	req := meta.NewListSpacesReq()

	resp, err := c.meta.ListSpaces(req)
	if err != nil {
		return nil, err
	}

	if err := codeErrorIfHappened(resp.Code, nil); err != nil {
		return nil, err
	}

	return newSpacesWrapper(resp.Spaces), nil
}

func (c *defaultMetaClient) Balance(req types.BalanceReq) (types.Balancer, error) {
	paras := make([][]byte, 0)

	var cmd meta.AdminCmd
	switch req.Cmd {
	case types.BalanceLeader:
		cmd = meta.AdminCmd_LEADER_BALANCE
	case types.BalanceData:
		cmd = meta.AdminCmd_DATA_BALANCE
	case types.BalanceDataRemove:
		cmd = meta.AdminCmd_DATA_BALANCE
		for _, ep := range req.HostsToRemove {
			paras = append(paras, []byte(ep))
		}
	default:
		return nil, nerrors.ErrUnsupported
	}

	paras = append(paras, []byte(req.Space))

	metaReq := &meta.AdminJobReq{
		Op:    meta.AdminJobOp_ADD,
		Cmd:   cmd,
		Paras: paras,
	}

	resp, err := c.meta.RunAdminJob(metaReq)
	if err != nil {
		return nil, err
	}

	if err := codeErrorIfHappened(resp.Code, nil); err != nil {
		return nil, err
	}

	return newBalancerWrap(c.meta, req.Space, resp.Result_.JobID), nil
}
