package v3_0

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	nthrift "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v3_0"
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

func (f *defaultFactoryDriver) NewValueBuilder() types.ValueBuilder {
	builder := nthrift.NewValueBuilder()
	return &valueBuilder{builder}
}

func (f *defaultFactoryDriver) NewDateBuilder() types.DateBuilder {
	builder := nthrift.NewDateBuilder()
	return &dateBuilder{builder}
}

func (f *defaultFactoryDriver) NewTimeBuilder() types.TimeBuilder {
	builder := nthrift.NewTimeBuilder()
	return &timeBuilder{builder}
}

func (f *defaultFactoryDriver) NewDateTimeBuilder() types.DateTimeBuilder {
	builder := nthrift.NewDateTimeBuilder()
	return &dateTimeBuilder{builder}
}

func (f *defaultFactoryDriver) NewEdgeBuilder() types.EdgeBuilder {
	builder := nthrift.NewEdgeBuilder()
	return &edgeBuilder{builder}
}

func (f *defaultFactoryDriver) NewNListBuilder() types.NListBuilder {
	builder := nthrift.NewNListBuilder()
	return &nListBuilder{builder}
}

func (f *defaultFactoryDriver) NewNMapBuilder() types.NMapBuilder {
	builder := nthrift.NewNMapBuilder()
	return &nMapBuilder{builder}
}
