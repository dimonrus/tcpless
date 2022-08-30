package tcpless

import (
	"net"
	"sync"
)

// concurrentClient client for concurrent request
type concurrentClient struct {
	m sync.RWMutex
	// concurrent clients
	clients []IClient
	// current client
	cursor int
	// max parallel requests
	concurrent int
	// buffer pool
	buffer *buffer
}

// Close all connection and release buffer
func (c *concurrentClient) Close() error {
	c.m.Lock()
	defer c.m.Unlock()
	var err error
	for i := range c.clients {
		err = c.clients[i].Close()
		c.buffer.Release(c.clients[i].Stream().Index())
	}
	return err
}

// set concurrent count
func (c *concurrentClient) concurrentCount(n int) {
	if n > 0 {
		c.concurrent = n
	} else {
		c.concurrent = 1
	}
}

// init buffers
func (c *concurrentClient) initBuffers(bufferSize int) {
	// init buffer
	if bufferSize == 0 {
		bufferSize = MinimumSharedBufferSize
	}
	c.buffer = CreateBuffer(c.concurrent, bufferSize)
}

// dial to server
func (c *concurrentClient) dialClients(address net.Addr, client func() IClient) error {
	c.clients = make([]IClient, c.concurrent)
	for i := 0; i < c.concurrent; i++ {
		buf, index := c.buffer.Pull()
		c.clients[i] = client()
		err := c.clients[i].Dial(address)
		if err != nil {
			if c.clients[i].Logger() != nil {
				c.clients[i].Logger().Errorln(err)
			}
			return err
		}
		c.clients[i].SetStream(newConnection(c.clients[i].Stream().Connection(), buf, index))
	}
	return nil
}

// ConcurrentClient create concurrent client with n, n <= 0 ignores
// bufferSize - shared buffer size
func ConcurrentClient(address net.Addr, n int, bufferSize int, client func() IClient) *concurrentClient {
	c := &concurrentClient{}
	// set n
	c.concurrentCount(n)
	// init buffers
	c.initBuffers(bufferSize)
	// construct clients
	_ = c.dialClients(address, client)
	return c
}
