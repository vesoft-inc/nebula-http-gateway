package client

import (
	"errors"

	"github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_0_0_ga"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_5_0"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_5_1"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/adapters/clientv2_6_0"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/logger"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/types"
)

func NewClient(addresses []types.HostAddress, config types.ClientConfig, logger logger.Logger) (types.Client, error) {
	switch config.Ver {
	case types.Version_2_0_0_ga:
		if config.SslConfig != nil {
			return nil, errors.New("ssl client not support")
		}
		return clientv2_0_0_ga.NewClient(addresses, config, logger)
	case types.Version_2_5_0:
		if config.SslConfig != nil {
			return nil, errors.New("ssl client not support")
		}
		return clientv2_5_0.NewClient(addresses, config, logger)
	case types.Version_2_5_1:
		if config.SslConfig != nil {
			return nil, errors.New("ssl client not support")
		}
		return clientv2_5_1.NewClient(addresses, config, logger)
	case types.Version_2_6_0:
		return clientv2_6_0.NewClient(addresses, config, logger)
	default:
		return nil, errors.New("client version not support")
	}
}
