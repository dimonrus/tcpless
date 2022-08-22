package tcpless

import (
	"github.com/dimonrus/gocli"
	"net"
	"time"
)

// Common parts for configs
type options struct {
	config Config
	logger gocli.Logger
}

// Config server configuration
type Config struct {
	// type of tcp network
	Network string
	// tcp address
	Address net.TCPAddr
	// connection limits
	Limits ConnectionLimit
}

// ConnectionLimit limits for connection
type ConnectionLimit struct {
	// Maximum number of connection
	MaxConnections uint16
	// Max idle time befor connection will be closed
	MaxIdle time.Duration
}
