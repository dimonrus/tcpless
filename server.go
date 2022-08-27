package tcpless

import (
	"github.com/dimonrus/gocli"
	"net"
)

// Server server main struct
type Server struct {
	options
	pool    *pool
	handler Handler
}

// Start server tcp connections
func (s *Server) Start() error {
	var err error
	s.pool.listener, err = net.Listen(s.config.Address.Network(), s.config.Address.String())
	go s.pool.idle()
	return err
}

// Stop close all tcp connections
func (s *Server) Stop() {

}

// Restart stop and start tcp serving
func (s *Server) Restart() {

}

// NewServer init new server
func NewServer(config Config, handler Handler, logger gocli.Logger) *Server {
	opt := options{
		config: config,
		logger: logger,
	}
	buffers := make([][]byte, config.Limits.MaxConnections)
	for i := uint16(0); i < config.Limits.MaxConnections; i++ {
		buffers[i] = make([]byte, config.Limits.SharedBufferSize)
	}
	return &Server{
		handler: handler,
		options: opt,
		pool: &pool{
			buffers:     buffers,
			options:     opt,
			connections: make([]*connection, 0),
		},
	}
}
