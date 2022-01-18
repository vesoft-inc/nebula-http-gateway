package nebula

import (
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

type (
	MetaClient interface {
		Open() error
		AddHosts(endpoints []string) error
		DropHosts(endpoints []string) error
		ListSpaces() (types.Spaces, error)
		BalanceData(space string) (types.Balancer, error)
		BalanceLeader(space string) (types.Balancer, error)
		BalanceDataRemove(space string, endpoints []string) (types.Balancer, error)
		Close() error
	}

	defaultMetaClient defaultClient
)

func NewMetaClient(endpoints []string, opts ...Option) (MetaClient, error) {
	c, err := NewClient(ConnectionInfo{
		MetaEndpoints: endpoints,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return c.Meta(), nil
}

func (c *defaultMetaClient) Open() error {
	return c.defaultClient().initDriver(func(driver types.Driver) error {
		return c.meta.open(driver)
	})
}

func (c *defaultMetaClient) Close() error {
	return c.meta.close()
}

func (c *defaultMetaClient) AddHosts(endpoints []string) error {
	return c.meta.AddHosts(endpoints)
}

func (c *defaultMetaClient) DropHosts(endpoints []string) error {
	return c.meta.DropHosts(endpoints)
}

func (c *defaultMetaClient) ListSpaces() (types.Spaces, error) {
	return c.meta.ListSpaces()
}

func (c *defaultMetaClient) BalanceData(space string) (types.Balancer, error) {
	return c.meta.Balance(types.BalanceReq{
		Cmd:   types.BalanceData,
		Space: space,
	})
}

func (c *defaultMetaClient) BalanceLeader(space string) (types.Balancer, error) {
	return c.meta.Balance(types.BalanceReq{
		Cmd:   types.BalanceLeader,
		Space: space,
	})
}

func (c *defaultMetaClient) BalanceDataRemove(space string, endpoints []string) (types.Balancer, error) {
	return c.meta.Balance(types.BalanceReq{
		Cmd:           types.BalanceDataRemove,
		Space:         space,
		HostsToRemove: endpoints,
	})
}

func (c *defaultMetaClient) defaultClient() *defaultClient {
	return (*defaultClient)(c)
}
