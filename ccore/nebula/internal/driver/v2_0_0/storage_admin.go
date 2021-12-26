package v2_0_0

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	storage2_0_0 "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_0_0/storage"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

var (
	_ types.StorageAdminClientDriver = (*defaultStorageAdminClient)(nil)
)

type (
	defaultStorageAdminClient struct {
		storageAdmin *storage2_0_0.StorageAdminServiceClient
	}
)

func newStorageAdminClient(transport thrift.Transport, pf thrift.ProtocolFactory) types.StorageAdminClientDriver {
	return &defaultStorageAdminClient{
		storageAdmin: storage2_0_0.NewStorageAdminServiceClientFactory(transport, pf),
	}
}

func (c *defaultStorageAdminClient) Open() error {
	return c.storageAdmin.Transport.Open()
}

func (c *defaultStorageAdminClient) Close() error {
	if c.storageAdmin != nil {
		if err := c.storageAdmin.Close(); err != nil {
			return err
		}
	}
	return nil
}
