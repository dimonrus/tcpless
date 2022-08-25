package tcpless

import (
	"context"
	"net"
	"sync"
	"time"
)

type pool struct {
	options
	m           sync.RWMutex
	connections []*connection
	listener    *net.TCPListener
}

func (p *pool) connection() *connection {
	for p.len() >= p.config.Limits.MaxConnections {
		time.Sleep(500 * time.Millisecond)
	}
	conn := p.accept()
	if conn != nil {
		p.m.Lock()
		p.connections = append(p.connections, conn)
		p.m.Unlock()
	}
	return conn
}

func (p *pool) len() uint16 {
	p.m.RLock()
	defer p.m.RUnlock()
	return uint16(len(p.connections))
}

func (p *pool) accept() *connection {
	conn, err := p.listener.AcceptTCP()
	if err != nil {
		p.logger.Errorln(err.Error())
		return nil
	}
	return &connection{
		TCPConn: conn,
		done:    make(chan struct{}),
	}
}

func (p *pool) removeConnection(c *connection) {
	p.m.Lock()
	defer p.m.Unlock()
	_ = c.Close()
	var i int
	for i = range p.connections {
		if p.connections[i] == c {
			break
		}
	}
	p.connections = append(p.connections[:i], p.connections[i+1:]...)
}

func (p *pool) process(c *connection) {
	for {
		select {
		case <-c.done:
			p.removeConnection(c)
			return
		default:
			sig := GobSignature{stream: c}
			err := sig.Decode(c.TCPConn)
			if err != nil {
				p.removeConnection(c)
				return
			} else {
				if callback, ok := registry[sig.route]; ok {
					callback(context.Background(), &sig)
				}
			}
		}
	}
}

func (p *pool) idle() {
	for {
		c := p.connection()
		go p.process(c)
	}
}
