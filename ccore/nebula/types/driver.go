package types

import (
	"sync"

	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	nerrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[Version]Driver)

	factoryDriversMu sync.RWMutex
	factoryDrivers   = make(map[Version]FactoryDriver)
)

type (
	Driver interface {
		NewGraphClientDriver(thrift.Transport, thrift.ProtocolFactory) GraphClientDriver
		NewMetaClientDriver(thrift.Transport, thrift.ProtocolFactory) MetaClientDriver
		NewStorageClientDriver(thrift.Transport, thrift.ProtocolFactory) StorageAdminClientDriver
	}

	GraphClientDriver interface {
		Open() error
		VerifyClientVersion() error
		Authenticate(username, password string) (AuthResponse, error)
		Signout(sessionId int64) (err error)
		Execute(sessionId int64, stmt []byte) (ExecutionResponse, error)
		ExecuteJson(sessionId int64, stmt []byte) ([]byte, error)
		Close() error
	}

	MetaClientDriver interface {
		Open() error
		VerifyClientVersion() error
		Close() error
	}

	StorageAdminClientDriver interface {
		Open() error
		Close() error
	}

	AuthResponse interface {
		SessionID() *int64
	}

	ExecutionResponse interface {
		GetLatencyInUs() int64
		GetData() DataSet
		GetSpaceName() []byte
		GetPlanDesc() PlanDescription
		GetComment() []byte
		IsSetData() bool
		IsSetSpaceName() bool
		IsSetErrorMsg() bool
		IsSetPlanDesc() bool
		IsSetComment() bool
		String() string
	}

	FactoryDriver interface {
		//NewDataset() DataSet
		//NewRow() Row
		NewValue() Value
		//NewDate() Date
		//NewTime() Time
		//NewDateTime() DateTime
		//NewVertex() Vertex
		//NewEdge() Edge
		//NewPath() Path
		//NewNList() NList
		//NewNMap() NMap
		//NewNSet() NSet
		//NewGeography() Geography
		//NewTag() Tag
		//NewStep() Step
		//NewPoint() Point
		//NewLineString() LineString
		//NewPolygon() Polygon
		//NewCoordinate() Coordinate
		//NewPlanDescription() PlanDescription
		//NewPlanNodeDescription() PlanNodeDescription
		//NewPair() Pair
		//NewProfilingStats() ProfilingStats
		//NewPlanNodeBranchInfo() PlanNodeBranchInfo
	}
)

func Register(version Version, driver Driver, factory FactoryDriver) {
	registerDriver(version, driver)
	registerFactoryDriver(version, factory)
}

func registerDriver(version Version, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("nebula: Register driver is nil")
	}
	if _, dup := drivers[version]; dup {
		panic("nebula: Register called twice for driver " + version)
	}
	drivers[version] = driver
}

func registerFactoryDriver(version Version, factory FactoryDriver) {
	factoryDriversMu.Lock()
	defer factoryDriversMu.Unlock()
	if factory == nil {
		panic("nebula: Register factory driver is nil")
	}
	if _, dup := factoryDrivers[version]; dup {
		panic("nebula: Register called twice for factory driver " + version)
	}
	factoryDrivers[version] = factory
}

func Drivers() []Version {
	driversMu.RLock()
	defer driversMu.RUnlock()
	list := make([]Version, 0, len(drivers))
	for version := range drivers {
		list = append(list, version)
	}
	return list
}

func GetDriver(version Version) (Driver, error) {
	driversMu.RLock()
	driver, ok := drivers[version]
	driversMu.RUnlock()
	if !ok {
		return nil, nerrors.ErrUnsupportedVersion
	}
	return driver, nil
}

func GetFactoryDriver(version Version) (FactoryDriver, error) {
	factoryDriversMu.RLock()
	factory, ok := factoryDrivers[version]
	factoryDriversMu.RUnlock()
	if !ok {
		return nil, nerrors.ErrUnsupportedVersion
	}
	return factory, nil
}
