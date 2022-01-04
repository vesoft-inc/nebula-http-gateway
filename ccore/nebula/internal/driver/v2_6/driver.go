package v2_6

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

func (f *defaultFactoryDriver) NewValue() types.Value {
	value := newValue()
	return newValueWrapper(value)
}

func (f *defaultFactoryDriver) NewDate() types.Date {
	date := newDate()
	return newDateWrapper(date)
}

func (f *defaultFactoryDriver) NewTime() types.Time {
	time := newTime()
	return newTimeWrapper(time)
}

func (f *defaultFactoryDriver) NewDateTime() types.DateTime {
	dateTime := newDateTime()
	return newDateTimeWrapper(dateTime)
}

func (f *defaultFactoryDriver) NewEdge() types.Edge {
	edge := newEdge()
	return newEdgeWrapper(edge)
}
