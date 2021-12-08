package clientv2_0_0_ga

import (
	"errors"
	"time"

	nebula "github.com/vesoft-inc/nebula-go/v2"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/client/logger"
	"github.com/vesoft-inc/nebula-http-gateway/pkg/client/types"
)

type Client struct {
	pool *nebula.ConnectionPool
	Ver  types.Version
	Ssl  bool
}

func NewClient(addresses []types.HostAddress, config types.ClientConfig, logger logger.Logger) (types.Client, error) {
	if config.SslConfig != nil {
		return nil, errors.New("ssl client not support")
	}
	pool, err := nebula.NewConnectionPool(hostsWrapper(addresses), configWrapper(config), logger)
	if err != nil {
		return nil, err
	}

	client := new(Client)
	client.Ver = config.Ver
	client.Ssl = false
	client.pool = pool

	return client, nil
}

func (c *Client) Close() {
	c.pool.Close()
}

func (c *Client) NewSession(account types.Account) (types.Session, error) {
	s, err := c.pool.GetSession(account.Username, account.Password)
	if err != nil {
		return nil, err
	}

	session := new(Session)
	session.session = s

	return session, nil
}

func (c *Client) Ping(address types.HostAddress, duration time.Duration) error {
	return c.pool.Ping(hostWrapper(address), duration)
}

func (c *Client) Version() types.Version {
	return c.Ver
}
