package v2_5_1

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	storage2_5_1 "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_5_1/storage"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

var (
	_ types.StorageAdminClientDriver = (*defaultStorageAdminClient)(nil)
)

type (
	defaultStorageAdminClient struct {
		storageAdmin *storage2_5_1.StorageAdminServiceClient
	}
)

func newStorageAdminClient(transport thrift.Transport, pf thrift.ProtocolFactory) types.StorageAdminClientDriver {
	return &defaultStorageAdminClient{
		storageAdmin: storage2_5_1.NewStorageAdminServiceClientFactory(transport, pf),
	}
}

func (c *defaultStorageAdminClient) Open() error {
	return c.storageAdmin.Open()
}

func (c *defaultStorageAdminClient) Close() error{
	if c.storageAdmin != nil {
		if err := c.storageAdmin.Close(); err != nil {
			return err
		}
	}
	return nil
}
