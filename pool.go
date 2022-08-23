package tcpless

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"
)

type pool struct {
	options
	m           sync.RWMutex
	connections []*connection
	cursor      uint16
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
	var connections = make([]*connection, len(p.connections)-1)
	var j int
	for i := range p.connections {
		if p.connections[i] == c {
			continue
		}
		connections[j] = p.connections[i]
		j++
	}
	p.connections = connections
	p.cursor = 0
}

func (p *pool) process(c *connection) {
	var nameLen [1]byte
	var dataLen [8]byte
	for {
		select {
		case <-c.done:
			p.removeConnection(c)
			return
		default:
			n, err := c.Read(nameLen[:])
			if err != nil {
				p.removeConnection(c)
				return
			} else {
				if n == 0 {
					continue
				}
				name := make([]byte, nameLen[0])
				n, err = c.Read(name)
				if err != nil {
					p.removeConnection(c)
					return
				}
				if callback, ok := registry[string(name)]; ok {
					n, err = c.Read(dataLen[:])
					if err != nil {
						p.removeConnection(c)
						return
					}
					sig := HandlerData{
						route:         string(name),
						contentLength: binary.BigEndian.Uint64(dataLen[:]),
						stream:        io.ReadWriteCloser(c.TCPConn),
					}
					callback(context.Background(), sig)
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
