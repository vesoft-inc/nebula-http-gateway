package v2_5

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_5/storage"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

var (
	_ types.StorageAdminClientDriver = (*defaultStorageAdminClient)(nil)
)

type (
	defaultStorageAdminClient struct {
		storageAdmin *storage.StorageAdminServiceClient
	}
)

func newStorageAdminClient(transport thrift.Transport, pf thrift.ProtocolFactory) types.StorageAdminClientDriver {
	return &defaultStorageAdminClient{
		storageAdmin: storage.NewStorageAdminServiceClientFactory(transport, pf),
	}
}

func (c *defaultStorageAdminClient) VerifyClientVersion() error {
	// v2.5 is not support verify client version, and it's the lowest version, so return not error.
	return nil
}

func (c *defaultStorageAdminClient) Open() error {
	return c.storageAdmin.Open()
}

func (c *defaultStorageAdminClient) Close() error {
	if c.storageAdmin != nil {
		if err := c.storageAdmin.Close(); err != nil {
			return err
		}
	}
	return nil
}
