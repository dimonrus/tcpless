package tcpless

import (
	"bytes"
	"net"
)

// Connection interface
type Connection interface {
	// Buffer get buffer
	Buffer() *bytes.Buffer
	// CleanBuffer clean buffer
	CleanBuffer() *bytes.Buffer
	// Connection get connection
	Connection() net.Conn
	// Index get index
	Index() uint16
	// Release connection
	Release() error
}

// Connection structure
type connection struct {
	net.Conn
	// is connection can be released
	done chan struct{}
	// all data should be stored here
	buffer *bytes.Buffer
	// index of buffer
	index uint16
}

// Buffer get buffer
func (c *connection) Buffer() *bytes.Buffer {
	return c.buffer
}

// CleanBuffer clean buffer
func (c *connection) CleanBuffer() *bytes.Buffer {
	c.buffer.Reset()
	return c.buffer
}

// Connection get connection
func (c *connection) Connection() net.Conn {
	return c.Conn
}

// Index get index buffer
func (c *connection) Index() uint16 {
	return c.index
}

// Release connection
func (c *connection) Release() error {
	c.done <- struct{}{}
	return c.Connection().Close()
}
