package types

import (
	"sync"

	"github.com/facebook/fbthrift/thrift/lib/go/thrift"
	nerrors "github.com/vesoft-inc/nebula-http-gateway/ccore/nebula/errors"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[Version]Driver)
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
)

func Register(version Version, driver Driver) {
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
