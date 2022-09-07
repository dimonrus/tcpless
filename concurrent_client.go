package tcpless

import (
	"github.com/dimonrus/gocli"
	"sync"
)

// concurrentClient GetFreeClient for concurrent request
type concurrentClient struct {
	m sync.RWMutex
	// concurrent clients
	clients []IClient
	// current GetFreeClient
	free []bool
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
	if c.free == nil {
		c.free = make([]bool, 0, c.concurrent)
	}
	for i := len(c.clients); i < c.concurrent; i++ {
		buf, index := c.buffer.Pull()
		client := constructor(config, logger)
		_, err := client.Dial()
		if err != nil {
			return err
		}
		client.SetStream(NewConnection(client.Stream().Connection(), buf, index))
		c.clients = append(c.clients, client)
		c.free = append(c.free, true)
	}
	return nil
}

// RegisterType register type in GetFreeClient
func (c *concurrentClient) RegisterType(v ...any) {
	c.m.RLock()
	defer c.m.RUnlock()
	for i := range c.clients {
		for _, t := range v {
			c.clients[i].RegisterType(t)
		}
	}
}

// GetConcurrent get number of concurrent connections
func (c *concurrentClient) GetConcurrent() (n int) {
	return c.concurrent
}

// GetFreeClient get free stream
func (c *concurrentClient) GetFreeClient() (client IClient, i int) {
	c.m.Lock()
	defer c.m.Unlock()
	for ; !c.free[i]; i++ {
		if i == c.concurrent {
			i = 0
		}
	}
	c.free[i] = false
	client = c.clients[i]
	return
}

// ReleaseClient release by index
func (c *concurrentClient) ReleaseClient(index int) {
	c.m.Lock()
	defer c.m.Unlock()
	c.free[index] = true
	return
}

// Call concurrent ask
// route - URI to server handler& Example: api.v1.hello
// request - chan with collection of requests
// processor - process handler if server respond to client
func (c *concurrentClient) Call(route string, request chan any, processor Handler) {
	for i := 0; i < c.concurrent; i++ {
		go func(route string, request chan any) {
			client, index := c.GetFreeClient()
			defer c.ReleaseClient(index)
			for r := range request {
				err := client.Ask(route, r)
				if err != nil {
					c.options.logger.Errorln(err)
					return
				}
				processor(client)
			}
		}(route, request)
	}
	return
}

// ConcurrentClient create concurrent GetFreeClient with n, n <= 0 ignores
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
