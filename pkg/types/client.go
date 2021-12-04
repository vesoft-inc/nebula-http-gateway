package types

import (
	"crypto/tls"
	"time"
)

type Client interface {
	Close()
	NewSession(account Account) (Session, error)
	Ping(address HostAddress, duration time.Duration) error
	Version() Version
}

type ClientConfig struct {
	Ver             Version
	TimeOut         time.Duration
	IdleTime        time.Duration
	MaxConnPoolSize int
	MinConnPoolSize int
	SslConfig       *tls.Config
}
