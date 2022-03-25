package nebula

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"

	nerrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

type (
	MetaClient interface {
		Open() error
		AddHosts(endpoints []string) (types.MetaBaser, error)
		DropHosts(endpoints []string) (types.MetaBaser, error)
		ListSpaces() (types.Spaces, error)
		BalanceData(space string) (types.Balancer, error)
		BalanceLeader(space string) (types.Balancer, error)
		BalanceDataRemove(space string, endpoints []string) (types.Balancer, error)
		ListHosts() (types.Hosts, error)
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
		return c.openRetry(driver)
	})
}

func (c *defaultMetaClient) ListHosts() (resp types.Hosts, err error) {
	retryErr := c.retryDo(func() (types.MetaBaser, error) {
		resp, err = c.meta.ListHosts()
		return resp, err
	})
	if retryErr != nil {
		return nil, retryErr
	}

	return
}

func (c *defaultMetaClient) Close() error {
	return c.meta.close()
}

func (c *defaultMetaClient) AddHosts(endpoints []string) (resp types.MetaBaser, err error) {
	retryErr := c.retryDo(func() (types.MetaBaser, error) {
		resp, err = c.meta.AddHosts(endpoints)
		return resp, err
	})
	if retryErr != nil {
		return nil, retryErr
	}

	return
}

func (c *defaultMetaClient) DropHosts(endpoints []string) (resp types.MetaBaser, err error) {
	retryErr := c.retryDo(func() (types.MetaBaser, error) {
		resp, err = c.meta.DropHosts(endpoints)
		return resp, err
	})
	if retryErr != nil {
		return nil, retryErr
	}

	return
}

func (c *defaultMetaClient) ListSpaces() (resp types.Spaces, err error) {
	retryErr := c.retryDo(func() (types.MetaBaser, error) {
		resp, err = c.meta.ListSpaces()
		return resp, err
	})
	if retryErr != nil {
		return nil, retryErr
	}

	return
}

func (c *defaultMetaClient) BalanceData(space string) (resp types.Balancer, err error) {
	retryErr := c.retryDo(func() (types.MetaBaser, error) {
		resp, err = c.meta.Balance(types.BalanceReq{
			Cmd:   types.BalanceData,
			Space: space,
		})
		return resp, err
	})
	if retryErr != nil {
		return nil, retryErr
	}

	return
}

func (c *defaultMetaClient) BalanceLeader(space string) (resp types.Balancer, err error) {
	retryErr := c.retryDo(func() (types.MetaBaser, error) {
		resp, err = c.meta.Balance(types.BalanceReq{
			Cmd:   types.BalanceLeader,
			Space: space,
		})
		return resp, err
	})
	if retryErr != nil {
		return nil, retryErr
	}

	return
}

func (c *defaultMetaClient) BalanceDataRemove(space string, endpoints []string) (resp types.Balancer, err error) {
	retryErr := c.retryDo(func() (types.MetaBaser, error) {
		resp, err = c.meta.Balance(types.BalanceReq{
			Cmd:           types.BalanceDataRemove,
			Space:         space,
			HostsToRemove: endpoints,
		})
		return resp, err
	})

	if retryErr != nil {
		return nil, retryErr
	}

	return
}

func (c *defaultMetaClient) defaultClient() *defaultClient {
	return (*defaultClient)(c)
}

func (c *defaultMetaClient) retryDo(fn func() (types.MetaBaser, error)) error {
	resp, err := fn()
	if err != nil {
		// check if transport exception
		if err = c.reconnect(err); err != nil {
			return err
		}
		resp, err = fn()
		if err != nil {
			return err
		}
	}
	// check if leader change
	if resp.GetCode() == nerrors.ErrorCode_E_LEADER_CHANGED {
		if err = c.updateLeader(resp.GetLeader()); err != nil {
			return err
		}
		if _, err = fn(); err != nil {
			return err
		}
	}
	return nil
}

func (c *defaultMetaClient) updateLeader(endpoint string) error {
	if err := c.meta.connection.SetEndpointIfExists(endpoint); err != nil {
		return err
	}
	return c.openRetry(c.driver)
}

func (c *defaultMetaClient) reconnect(err error) error {
	if _, ok := err.(thrift.TransportException); !ok {
		return err
	}
	return c.openRetry(c.driver)
}

func (c *defaultMetaClient) openRetry(driver types.Driver) error {
	n := c.meta.connection.GetEndpointsLen()
	for i := 0; i < n; i++ {
		_ = c.meta.close()

		err := c.meta.open(driver)
		c.meta.connection.UpdateNextIndex() // update nextIndex every time
		if err == nil {
			return nil
		}
	}
	return nerrors.ErrNoValidMetaEndpoint
}
