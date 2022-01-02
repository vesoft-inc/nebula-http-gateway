package nebula

import (
	"sync"

	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	nerrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

type (
	// TODO: add client pool management
	// TODO: add version auto recognize

	Client interface {
		Graph() GraphClient
		Meta() MetaClient
		StorageAdmin() StorageAdminClient
		Factory() Factory
		Version() Version
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
		o              Options
		driver         types.Driver
		initDriverOnce sync.Once
		graph          *driverGraph
		meta           *driverMeta
		storageAdmin   *driverStorageAdmin
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

	return &defaultClient{
		o:            o,
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

func (c *defaultClient) Factory() Factory {
	f, _ := NewFactory(WithVersion(c.o.version))
	return f
}

func (c *defaultClient) Version() Version {
	return c.o.version
}

func (c *defaultClient) initDriver(checkFn func(types.Driver) error) error {
	if c.o.version != VersionAuto {
		driver, err := types.GetDriver(c.o.version)
		if err != nil {
			return err
		}
		if err = checkFn(driver); err != nil {
			return err
		}
		c.driver = driver
		return nil
	}

	for _, v := range Versions {
		driver, err := types.GetDriver(v)
		if err != nil {
			return err
		}
		if err = checkFn(driver); err != nil {
			if nerrors.IsCodeError(err, nerrors.ErrorCode_E_CLIENT_SERVER_INCOMPATIBLE) {
				continue
			}
			if e, ok := err.(thrift.ApplicationException); ok && e.TypeID() == thrift.WRONG_METHOD_NAME {
				continue
			}
			return err
		}
		c.driver = driver
		c.o.version = v
		return nil
	}
	return nerrors.ErrUnsupportedVersion
}
