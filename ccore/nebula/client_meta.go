package nebula

type (
	MetaClient interface {
		Open() error
		Close() error
	}

	defaultMetaClient defaultClient
)

func NewMetaClient(endpoints []string,  opts ...Option) (MetaClient, error) {
	c, err := NewClient(ConnectionInfo{
		MetaEndpoints: endpoints,
	}, opts...)
	if err != nil {
		return nil, err
	}
	return c.Meta(), nil
}

func (c *defaultMetaClient) Open() error {
	return c.meta.close()
}

func (c *defaultMetaClient) Close() error {
	return c.meta.open(c.driver)
}
