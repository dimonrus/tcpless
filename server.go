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
	client  ClientConstructor
}

// Start server tcp connections
func (s *Server) Start() error {
	var err error
	s.pool.listener, err = net.Listen(s.config.Address.Network(), s.config.Address.String())
	go s.idle()
	return err
}

// Stop close all tcp connections
func (s *Server) Stop() {

}

// Restart stop and start tcp serving
func (s *Server) Restart() {

}

// Idle listen connection
func (s *Server) idle() {
	for {
		c := s.pool.connection()
		client := s.client(s.logger)
		client.SetStream(c)
		go s.pool.process(client)
	}
}

// NewServer init new server
func NewServer(config Config, handler Handler, client ClientConstructor, logger gocli.Logger) *Server {
	opt := options{
		config: config,
		logger: logger,
	}
	return &Server{
		client:  client,
		handler: handler,
		options: opt,
		pool: &pool{
			buffer:      CreateBuffer(int(config.Limits.MaxConnections), int(config.Limits.SharedBufferSize)),
			options:     opt,
			connections: make([]*connection, 0),
		},
	}
}
