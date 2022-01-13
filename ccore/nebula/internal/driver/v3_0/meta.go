package v3_0

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	nthrift "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v3_0"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v3_0/meta"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
	"net"
	"strconv"
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
	return codeErrorIfHappened(resp.Code, []byte(nthrift.ErrorCodeToName[resp.Code]))
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
	return codeErrorIfHappened(resp.Code, []byte(nthrift.ErrorCodeToName[resp.Code]))
}

func (c *defaultMetaClient) ListSpaces() (types.ListSpacesResponse, error) {
	req := meta.NewListSpacesReq()

	resp, err := c.meta.ListSpaces(req)
	if err != nil {
		return nil, err
	}

	return newListSpacesResponseWrapper(resp.Spaces), codeErrorIfHappened(resp.Code, []byte(nthrift.ErrorCodeToName[resp.Code]))
}
