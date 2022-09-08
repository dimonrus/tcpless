package tcpless

import (
	"context"
	"net"
	"sync"
	"time"
)

type pool struct {
	options
	buffer      *buffer
	m           sync.RWMutex
	connections []*connection
	listener    net.Listener
}

func (p *pool) connection() *connection {
	conn := p.accept()
	if conn != nil {
		for p.len() >= p.config.Limits.MaxConnections {
			time.Sleep(500 * time.Millisecond)
		}
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
	conn, err := p.listener.Accept()
	if err != nil {
		p.logger.Errorln(err.Error())
		return nil
	}
	c := &connection{
		Conn: conn,
		done: make(chan struct{}),
	}
	c.buffer, c.index = p.buffer.Pull()
	return c
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
	p.buffer.Release(c.index)
}

// release connection and exit from idle
func (p *pool) release() {
	p.m.Lock()
	defer p.m.Unlock()
	for _, c := range p.connections {
		c.done <- struct{}{}
	}
	_ = p.listener.Close()
	p.listener = nil
}

func (p *pool) process(client IClient) {
	for {
		select {
		case <-client.Stream().Exit():
			p.removeConnection(client.Stream().(*connection))
			return
		default:
			// read data into signature
			sig, err := client.Read()
			if err != nil {
				p.removeConnection(client.Stream().(*connection))
				return
			} else {
				if callback, ok := registry[sig.Route()]; ok {
					client.WithContext(context.Background())
					for _, hook := range routeRegistry.GetHooks(sig.Route()) {
						hook(client)
					}
					callback(client)
					// clear buffer after, reuse memory
					client.Stream().Buffer().Reset()
					sig.Reset()
				} else {
					p.removeConnection(client.Stream().(*connection))
					return
				}
			}
		}
	}
}
