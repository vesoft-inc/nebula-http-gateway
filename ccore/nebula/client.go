package nebula

import "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"

type (
	// TODO: add client pool management
	// TODO: add version auto recognize

	Client interface {
		Graph() GraphClient
		Meta() MetaClient
		StorageAdmin() StorageAdminClient
	}

	ConnectionInfo struct {
		GraphEndpoints        []string
		MetaEndpoints         []string
		StorageAdminEndpoints []string
		GraphAccount          Account
	}

	Account struct {
		Username string
		Password string
	}

	defaultClient struct {
		o            Options
		driver       types.Driver
		graph        *driverGraph
		meta         *driverMeta
		storageAdmin *driverStorageAdmin
	}
)

func NewClient(info ConnectionInfo, opts ...Option) (Client, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}
	o.complete()
	if err := o.validate(); err != nil {
		return nil, err
	}

	driver, err := types.GetDriver(o.version)
	if err != nil {
		return nil, err
	}

	return &defaultClient{
		o:            o,
		driver:       driver,
		graph:        newDriverGraph(info.GraphEndpoints, info.GraphAccount.Username, info.GraphAccount.Password, &o.graph),
		meta:         newDriverMeta(info.MetaEndpoints, &o.meta),
		storageAdmin: newDriverStorageAdmin(info.StorageAdminEndpoints, &o.storageAdmin),
	}, nil
}

func (c *defaultClient) Graph() GraphClient {
	return (*defaultGraphClient)(c)
}

func (c *defaultClient) Meta() MetaClient {
	return (*defaultMetaClient)(c)
}

func (c *defaultClient) StorageAdmin() StorageAdminClient {
	return (*defaultStorageAdminClient)(c)
}
