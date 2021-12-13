package nebula

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
	return c.storageAdmin.open(c.driver)
}

func (c *defaultStorageAdminClient) Close() error {
	return c.storageAdmin.close()
}
