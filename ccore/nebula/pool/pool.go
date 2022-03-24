package pool

import (
	"container/list"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vesoft-inc/nebula-http-gateway/ccore/nebula"
)

type Pooler interface {
	Get() (*Conn, error)
	Put(*Conn)
	Remove(*Conn)
	Replace(*Conn) (*Conn, error)

	Close()
}

type Options struct {
	TimeOut         time.Duration
	IdleTimeOut     time.Duration
	MaxConnPoolSize int
	MinConnPoolSize int
}

func (o *Options) complete() {
	if o.TimeOut <= 0 {
		o.TimeOut = 3 * time.Second
	}
	if o.IdleTimeOut < 0 {
		o.IdleTimeOut = 30 * time.Minute
	}
	if o.MaxConnPoolSize <= 0 {
		o.MaxConnPoolSize = 10
	}
	if o.MinConnPoolSize < 0 {
		o.MinConnPoolSize = 0
	}
}

type GraphClientPool struct {
	opt *Options

	endpoints     []string
	endpointIndex int
	username      string
	password      string

	idleConns   list.List
	activeConns list.List
	connsMu     sync.Mutex

	closedCh chan struct{}
	closed   bool
}

var _ Pooler = (*GraphClientPool)(nil)

func NewGraphClientPool(endpoints []string, username, password string, opt *Options) (*GraphClientPool, error) {
	convEps, err := DomainToIP(endpoints)
	if err != nil {
		return nil, err
	}

	opt.complete()

	p := &GraphClientPool{
		opt: opt,

		endpoints:     convEps,
		endpointIndex: 0,
		username:      username,
		password:      password,
		closedCh:      make(chan struct{}),
	}

	if err = p.checkMinIdleConns(); err != nil {
		return nil, err
	}

	if opt.IdleTimeOut > 0 {
		go p.reaper()
	}

	return p, nil
}

func (p *GraphClientPool) Get() (*Conn, error) {
	p.connsMu.Lock()
	defer p.connsMu.Unlock()
	return p.getIdleConn()
}

func (p *GraphClientPool) Put(cn *Conn) {
	p.connsMu.Lock()
	defer p.connsMu.Unlock()
	p.putConn(cn)
}

func (p *GraphClientPool) Remove(cn *Conn) {
	p.connsMu.Lock()
	defer p.connsMu.Unlock()
	p.removeConn(cn)
}

func (p *GraphClientPool) Replace(cn *Conn) (*Conn, error) {
	p.connsMu.Lock()
	defer p.connsMu.Unlock()
	p.removeConn(cn)
	return p.getIdleConn()
}

func (p *GraphClientPool) Close() {
	p.connsMu.Lock()
	defer p.connsMu.Unlock()
	idleLen := p.idleConns.Len()
	activeLen := p.activeConns.Len()

	for i := 0; i < idleLen; i++ {
		_ = p.idleConns.Front().Value.(*Conn).Close()
		p.idleConns.Remove(p.idleConns.Front())
	}
	for i := 0; i < activeLen; i++ {
		_ = p.activeConns.Front().Value.(*Conn).Close()
		p.activeConns.Remove(p.activeConns.Front())
	}

	p.closed = true
	if p.closedCh != nil {
		close(p.closedCh)
	}
}

func (p *GraphClientPool) checkMinIdleConns() error {
	for i := 0; i < p.opt.MinConnPoolSize; i++ {
		newConn, err := NewConn(p.endpoints[i%len(p.endpoints)], p.username, p.password,
			nebula.WithGraphTimeout(p.opt.TimeOut))
		if err != nil {
			idleLen := p.idleConns.Len()
			for j := 0; j < idleLen; j++ {
				_ = p.idleConns.Front().Value.(*Conn).Close()
				p.idleConns.Remove(p.idleConns.Front())
			}
			return fmt.Errorf("failed to open connection, error: %s ", err.Error())
		}
		p.idleConns.PushBack(newConn)
	}
	return nil
}

func (p *GraphClientPool) reaper() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if p.closed {
				return
			}
			staleConns := p.reapStaleConns()
			for _, c := range staleConns {
				_ = c.Close()
			}
		case <-p.closedCh:
			return
		}
	}
}

func (p *GraphClientPool) reapStaleConns() []*Conn {
	expiredSince := time.Now().Add(-p.opt.IdleTimeOut)
	var newEle *list.Element = nil

	maxCleanSize := p.idleConns.Len() + p.activeConns.Len() - p.opt.MinConnPoolSize

	staleConns := make([]*Conn, 0)
	for ele := p.idleConns.Front(); ele != nil; {
		if maxCleanSize == 0 {
			return staleConns
		}

		newEle = ele.Next()
		// Check connection is expired
		if !ele.Value.(*Conn).createdAt.Before(expiredSince) {
			return staleConns
		}
		staleConns = append(staleConns, ele.Value.(*Conn))
		p.idleConns.Remove(ele)
		ele = newEle
		maxCleanSize--
	}

	return staleConns
}

func (p *GraphClientPool) newConn() (*Conn, error) {
	totalConn := p.idleConns.Len() + p.activeConns.Len()

	if totalConn >= p.opt.MaxConnPoolSize {
		return nil, fmt.Errorf("failed to get connection: No valid connection" +
			" in the idle queue and connection number has reached the pool capacity")
	}

	newConn, err := NewConn(p.endpoints[p.endpointIndex%len(p.endpoints)], p.username, p.password,
		nebula.WithGraphTimeout(p.opt.TimeOut))
	p.endpointIndex++
	if err != nil {
		return nil, err
	}

	p.activeConns.PushBack(newConn)

	return newConn, err
}

func (p *GraphClientPool) getIdleConn() (*Conn, error) {
	if p.idleConns.Len() > 0 {
		newConn := p.idleConns.Front().Value.(*Conn)

		p.idleConns.Remove(p.idleConns.Front())
		p.activeConns.PushBack(newConn)
		return newConn, nil
	}

	// Create a new connection if there is no idle connection and total connection < pool max size
	newConn, err := p.newConn()
	return newConn, err
}

func (p *GraphClientPool) putConn(cn *Conn) {
	for ele := p.activeConns.Front(); ele != nil; ele = ele.Next() {
		if ele.Value.(*Conn) == cn {
			p.activeConns.Remove(ele)
			break
		}
	}
	p.idleConns.PushBack(cn)
}

func (p *GraphClientPool) removeConn(cn *Conn) {
	for ele := p.activeConns.Front(); ele != nil; ele = ele.Next() {
		if ele.Value.(*Conn) == cn {
			p.activeConns.Remove(ele)
			break
		}
	}
	_ = cn.Close()
}

func DomainToIP(endpoints []string) ([]string, error) {
	convEps := make([]string, 0)
	for _, ep := range endpoints {
		hostPort := strings.Split(ep, ":")
		if len(hostPort) != 2 {
			return nil, fmt.Errorf("invalid endpoint: %s", ep)
		}
		port, err := strconv.Atoi(hostPort[1])
		if err != nil {
			return nil, err
		}
		ip, err := net.LookupIP(hostPort[0])
		if err != nil {
			return nil, fmt.Errorf("could not get %s ip: %v\n", ep, err)
		}
		convEp := fmt.Sprintf("%s:%d", ip, port)
		convEps = append(convEps, convEp)
	}

	return convEps, nil
}
