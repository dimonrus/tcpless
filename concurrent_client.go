package tcpless

import (
	"github.com/dimonrus/gocli"
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
	// options
	options options
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
func (c *concurrentClient) dialClients(client ClientConstructor, config *Config, logger gocli.Logger) error {
	c.clients = make([]IClient, c.concurrent)
	for i := 0; i < c.concurrent; i++ {
		buf, index := c.buffer.Pull()
		c.clients[i] = client(config, logger)
		_, err := c.clients[i].Dial()
		if err != nil {
			if logger != nil {
				logger.Errorln(err)
			}
			return err
		}
		c.clients[i].SetStream(newConnection(c.clients[i].Stream().Connection(), buf, index))
	}
	return nil
}

// ConcurrentClient create concurrent client with n, n <= 0 ignores
// bufferSize - shared buffer size
func ConcurrentClient(n int, bufferSize int, client ClientConstructor, config *Config, logger gocli.Logger) *concurrentClient {
	c := &concurrentClient{}
	// set n
	c.concurrentCount(n)
	// init buffers
	c.initBuffers(bufferSize)
	// construct clients
	_ = c.dialClients(client, config, logger)
	return c
}
