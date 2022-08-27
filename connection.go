package tcpless

import (
	"net"
)

// Connection interface
type Connection interface {
	// Buffer get buffer
	Buffer() *PermanentBuffer
	// Connection get connection
	Connection() net.Conn
	// Exit chan for exit
	Exit() <-chan struct{}
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
	buffer *PermanentBuffer
	// index of buffer
	index uint16
}

// Buffer get buffer
func (c *connection) Buffer() *PermanentBuffer {
	return c.buffer
}

// Connection get connection
func (c *connection) Connection() net.Conn {
	return c.Conn
}

// Exit listen for exit
func (c *connection) Exit() <-chan struct{} {
	return c.done
}

// Index get index buffer
func (c *connection) Index() uint16 {
	return c.index
}

// Release connection
func (c *connection) Release() error {
	return c.Conn.Close()
}
