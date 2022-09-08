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
	// config
	config *Config
	// logger
	logger gocli.Logger
}

// Config server configuration
type Config struct {
	// tcp address
	Address *net.TCPAddr
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
	MaxConnections uint16 `yaml:"maxConnections"`
	// Max idle time before connection will be closed
	MaxIdle time.Duration `yaml:"maxIdle"`
	// Max process body size
	SharedBufferSize int32 `yaml:"sharedBufferSize"`
}
