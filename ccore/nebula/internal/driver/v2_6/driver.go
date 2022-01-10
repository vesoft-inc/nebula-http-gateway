package v2_6

import (
	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	nthrift "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/internal/thrift/v2_6"
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
	types.Register(types.Version2_6, &defaultDriver{}, &defaultFactoryDriver{})
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
	value := nthrift.NewValue()
	return &valueBuilder{value}
}

func (f *defaultFactoryDriver) NewDateBuilder() types.DateBuilder {
	date := nthrift.NewDate()
	return &dateBuilder{date}
}

func (f *defaultFactoryDriver) NewTimeBuilder() types.TimeBuilder {
	time := nthrift.NewTime()
	return &timeBuilder{time}
}

func (f *defaultFactoryDriver) NewDateTimeBuilder() types.DateTimeBuilder {
	dateTime := nthrift.NewDateTime()
	return &dateTimeBuilder{dateTime}
}

func (f *defaultFactoryDriver) NewEdgeBuilder() types.EdgeBuilder {
	edge := nthrift.NewEdge()
	return &edgeBuilder{edge}
}

func (f *defaultFactoryDriver) NewNListBuilder() types.NListBuilder {
	nlist := nthrift.NewNList()
	return &nListBuilder{nlist}
}

func (f *defaultFactoryDriver) NewNMapBuilder() types.NMapBuilder {
	nmap := nthrift.NewNMap()
	return &nMapBuilder{nmap}
}
