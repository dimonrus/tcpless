package tcpless

import "net"

// Connection structure
type connection struct {
	*net.TCPConn
	// is connection can be released
	done chan struct{}
}
