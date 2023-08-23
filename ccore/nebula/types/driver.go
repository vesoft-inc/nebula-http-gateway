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
		ExecuteWithParameter(sessionId int64, stmt []byte, params map[string]Value) (ExecutionResponse, error)
		Close() error
	}

	MetaClientDriver interface {
		Open() error
		VerifyClientVersion() error
		AddHosts(endpoints []string) (MetaBaser, error)
		DropHosts(endpoints []string) (MetaBaser, error)
		ListSpaces() (Spaces, error)
		Balance(req BalanceReq) (Balancer, error)
		ListHosts() (Hosts, error)
		ListZones() (Zones, error)
		Close() error
	}

	StorageAdminClientDriver interface {
		Open() error
		Close() error
	}

	AuthResponse interface {
		SessionID() *int64
		GetTimezoneInfo() TimezoneInfo
	}

	ExecutionResponse interface {
		GetLatencyInUs() int64
		GetData() DataSet
		GetSpaceName() []byte
		GetPlanDesc() PlanDescription
		GetComment() []byte
		GetErrorCode() nerrors.ErrorCode
		GetErrorMsg() []byte
		IsSetData() bool
		IsSetSpaceName() bool
		IsSetErrorMsg() bool
		IsSetPlanDesc() bool
		IsSetComment() bool
		String() string
	}

	Space interface {
		GetName() string
		GetId() int32
	}

	Spaces interface {
		MetaBaser
		GetSpaces() []Space
	}

	Balancer interface {
		MetaBaser
		GetStats() (BalanceStats, error)
	}

	FactoryDriver interface {
		NewValueBuilder() ValueBuilder
		NewDateBuilder() DateBuilder
		NewTimeBuilder() TimeBuilder
		NewDateTimeBuilder() DateTimeBuilder
		NewEdgeBuilder() EdgeBuilder
		NewNListBuilder() NListBuilder
		NewNMapBuilder() NMapBuilder
	}

	Host interface {
		GetHostItem() HostItem
	}

	Hosts interface {
		MetaBaser
		GetHosts() []Host
	}

	Zone interface {
		GetName() string
		GetHosts() []*HostAddr
	}

	Zones interface {
		MetaBaser
		GetZones() []Zone
	}

	Coder interface {
		GetCode() nerrors.ErrorCode
	}

	MetaBaser interface {
		Coder
		GetLeader() string
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
