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
	// ClientModeResponder Client can't ask
	ClientModeResponder = 0
	// ClientModeAsker client can't respond
	ClientModeAsker = 1
)

// Common parts for configs
type options struct {
	// Config
	config *Config
	// Logger
	logger gocli.Logger
}

// Config server configuration
type Config struct {
	// tcp address
	Address net.Addr
	// connection limits
	Limits ConnectionLimit
	// TLS configuration
	TLS TLSConfig
	// Mode client mode. 0 - Responder, - 1 Asker
	Mode uint8
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
