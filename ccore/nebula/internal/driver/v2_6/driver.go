package v2_6

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/types"
)

var (
	_ types.Driver = (*defaultDriver)(nil)
)

type (
	defaultDriver struct{}
)

func init() {
	types.Register(types.Version2_6, &defaultDriver{})
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
