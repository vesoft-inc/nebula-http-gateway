package clientv2_6_0

import (
	"errors"
	nebula "github.com/vesoft-inc/nebula-go/v2"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/client/types"
)

func configWrapper(config types.ClientConfig) nebula.PoolConfig {
	return nebula.PoolConfig{
		TimeOut:         config.TimeOut,
		IdleTime:        config.IdleTime,
		MaxConnPoolSize: config.MaxConnPoolSize,
		MinConnPoolSize: config.MinConnPoolSize,
	}
}

func hostWrapper(address types.HostAddress) nebula.HostAddress {
	return nebula.HostAddress{
		Host: address.Host,
		Port: address.Port,
	}
}

func hostsWrapper(hostAddresses []types.HostAddress) []nebula.HostAddress {
	hosts := make([]nebula.HostAddress, 0, len(hostAddresses))
	for i := range hostAddresses {
		hosts = append(hosts, hostWrapper(hostAddresses[i]))
	}
	return hosts
}

func ResultSetWrapper(rset types.ResultSet) (*nebula.ResultSet, error) {
	switch rset.(type) {
	case *nebula.ResultSet:
		return rset.(*nebula.ResultSet), nil
	default:
		return nil, errors.New("type not match")
	}
}
