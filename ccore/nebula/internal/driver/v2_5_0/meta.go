package v2_5_0

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	meta2_5_0 "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_5_0/meta"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

var (
	_ types.MetaClientDriver = (*defaultMetaClient)(nil)
)

type (
	defaultMetaClient struct {
		meta *meta2_5_0.MetaServiceClient
	}
)

func newMetaClient(transport thrift.Transport, pf thrift.ProtocolFactory) types.MetaClientDriver {
	return &defaultMetaClient{
		meta: meta2_5_0.NewMetaServiceClientFactory(transport, pf),
	}
}

func (c *defaultMetaClient) Open() error {
	return c.meta.Open()
}

func (c *defaultMetaClient) Close() error {
	if c.meta != nil {
		if err := c.meta.Close(); err != nil {
			return err
		}
	}
	return nil
}
