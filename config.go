package tcpless

import (
	"github.com/dimonrus/gocli"
	"net"
	"time"
)

const (
	// DefaultSharedBufferSize - 4 MB
	DefaultSharedBufferSize = 1024 * 1024 * 4 // 4 MB
	// MinimumSharedBufferSize - 2 KB
	MinimumSharedBufferSize = 1024 * 2 // 2 KB
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
