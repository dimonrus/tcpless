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
func (c *concurrentClient) dialClients(constructor ClientConstructor, config *Config, logger gocli.Logger) error {
	c.m.Lock()
	defer c.m.Unlock()
	if c.clients == nil {
		c.clients = make([]IClient, 0, c.concurrent)
	}
	for i := len(c.clients); i < c.concurrent; i++ {
		buf, index := c.buffer.Pull()
		client := constructor(config, logger)
		_, err := client.Dial()
		if err != nil {
			return err
		}
		client.SetStream(newConnection(client.Stream().Connection(), buf, index))
		c.clients = append(c.clients, client)
	}
	return nil
}

// RegisterType register type in client
func (c *concurrentClient) RegisterType(v ...any) {
	c.m.RLock()
	defer c.m.RUnlock()
	for i := range c.clients {
		for _, t := range v {
			c.clients[i].RegisterType(t)
		}
	}
}

// get current client
func (c *concurrentClient) client() (client IClient) {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.clients[c.cursor]
}

// Ask any data
func (c *concurrentClient) Ask(route string, request any, response any) (err error) {
	client := c.client()
	err = client.Ask(route, request)
	if err != nil || response == nil {
		return err
	}
	err = client.Parse(response)
	return
}

// ConcurrentClient create concurrent client with n, n <= 0 ignores
// bufferSize - shared buffer size
func ConcurrentClient(n int, bufferSize int, client ClientConstructor, config *Config, logger gocli.Logger) (*concurrentClient, error) {
	c := &concurrentClient{}
	// set n
	c.concurrentCount(n)
	// init buffers
	c.initBuffers(bufferSize)
	// construct clients
	return c, c.dialClients(client, config, logger)
}
