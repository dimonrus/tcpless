package tcpless

import (
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"
)

type pool struct {
	buffers [][]byte
	options
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
	return &connection{
		Conn: conn,
		done: make(chan struct{}),
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
	client := &GobClient{}
	client.SetStream(c)
	for {
		select {
		case <-c.done:
			p.removeConnection(c)
			return
		default:
			time.Sleep(time.Second)
			sig := client.Read()
			_ = sig
			//if sig == nil {
			//	p.removeConnection(c)
			//	return
			//} else {
			//	if callback, ok := registry[sig.Route()]; ok {
			//		callback(context.Background(), client, sig)
			//	}
			//}
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			report := make(map[string]string)
			report["allocated"] = fmt.Sprintf("%v KB", m.Alloc/1024)
			report["total_allocated"] = fmt.Sprintf("%v KB", m.TotalAlloc/1024)
			report["system"] = fmt.Sprintf("%v KB", m.Sys/1024)
			report["garbage_collectors"] = fmt.Sprintf("%v", m.NumGC)
			fmt.Println(report)
		}
	}
}

func (p *pool) idle() {
	for {
		c := p.connection()
		go p.process(c)
	}
}
