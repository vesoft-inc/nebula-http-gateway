package v3_0

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

var (
	_ types.Driver        = (*defaultDriver)(nil)
	_ types.FactoryDriver = (*defaultFactoryDriver)(nil)
)

type (
	defaultDriver        struct{}
	defaultFactoryDriver struct{}
)

func init() {
	types.Register(types.Version3_0, &defaultDriver{}, &defaultFactoryDriver{})
}

func (d *defaultDriver) NewGraphClientDriver(transport thrift.Transport, pf thrift.ProtocolFactory) types.GraphClientDriver {
	return newGraphClient(transport, pf)
}

func (d *defaultDriver) NewMetaClientDriver(transport thrift.Transport, pf thrift.ProtocolFactory) types.MetaClientDriver {
	return newMetaClient(transport, pf)
}

func (d *defaultDriver) NewStorageClientDriver(transport thrift.Transport, pf thrift.ProtocolFactory) types.StorageAdminClientDriver {
	return newStorageAdminClient(transport, pf)
}

func (f *defaultFactoryDriver) NewValue() types.Value {
	value := newValue()
	return newValueWrapper(value)
}
