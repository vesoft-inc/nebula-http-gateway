package nebula

import "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"

type (
	StorageAdminClient interface {
		Open() error
		Close() error
	}

	defaultStorageAdminClient defaultClient
)

func NewStorageAdminClient(endpoints []string, opts ...Option) (StorageAdminClient, error) {
	c, err := NewClient(ConnectionInfo{
		StorageAdminEndpoints: endpoints,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return c.StorageAdmin(), nil
}

func (c *defaultStorageAdminClient) Open() error {
	return c.defaultClient().initDriver(func(driver types.Driver) error {
		return c.storageAdmin.open(driver)
	})
}

func (c *defaultStorageAdminClient) Close() error {
	return c.storageAdmin.close()
}

func (c *defaultStorageAdminClient) defaultClient() *defaultClient {
	return (*defaultClient)(c)
}
