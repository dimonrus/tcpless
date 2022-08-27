package tcpless

import (
	"github.com/dimonrus/gocli"
	"net"
	"time"
)

const (
	DefaultSharedBufferSize = 1024 * 1024 * 4 // 4 MB
)

// Common parts for configs
type options struct {
	// Config
	config Config
	// Logger
	logger gocli.Logger
}

// Config server configuration
type Config struct {
	// tcp address
	Address net.Addr
	// connection limits
	Limits ConnectionLimit
}

// ConnectionLimit limits for connection
type ConnectionLimit struct {
	// Maximum number of connection
	MaxConnections uint16
	// Max idle time before connection will be closed
	MaxIdle time.Duration
	// Max process body size
	SharedBufferSize int32
}
