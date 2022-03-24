package pool

import (
	"time"

	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula"
)

type Conn struct {
	nebula.GraphClient
	createdAt time.Time
}

func NewConn(endpoint, username, pwd string, opts ...nebula.Option) (*Conn, error) {
	c, err := nebula.NewGraphClient([]string{endpoint}, username, pwd, opts...)
	if err != nil {
		return nil, err
	}

	if err = c.Open(); err != nil {
		return nil, err
	}

	return &Conn{
		GraphClient: c,
		createdAt:   time.Now(),
	}, nil
}
