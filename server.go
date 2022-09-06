package tcpless

import (
	"crypto/tls"
	"errors"
	"github.com/dimonrus/gocli"
	"net"
	"runtime"
)

// Server server main struct
type Server struct {
	options
	pool    *pool
	handler Handler
	client  ClientConstructor
	stop    chan struct{}
}

var tlsCertificateError = errors.New("no certificate found. Config is empty")

// Start server tcp connections
func (s *Server) Start() error {
	var err error
	if s.config.TLS.Enabled {
		var config *tls.Config
		config, err = s.config.TLS.LoadTLSConfig()
		if err != nil {
			return err
		}
		if config == nil {
			return tlsCertificateError
		}
		s.pool.m.Lock()
		s.pool.listener, err = tls.Listen(s.config.Address.Network(), s.config.Address.String(), config)
		s.pool.m.Unlock()
	} else {
		s.pool.m.Lock()
		s.pool.listener, err = net.Listen(s.config.Address.Network(), s.config.Address.String())
		s.pool.m.Unlock()
	}
	go s.idle()
	return err
}

// Stop close all tcp connections
func (s *Server) Stop() {
	s.stop <- struct{}{}
}

// Idle listen connection
func (s *Server) idle() {
	chCon := make(chan *connection)
	for {
		go func() {
			chCon <- s.pool.connection()
		}()
		select {
		// release all connections and stop pool process
		case <-s.stop:
			s.pool.release()

			runtime.GC()
			return
		// listen connections
		case c := <-chCon:
			client := s.client(s.config, s.logger)
			client.SetStream(c)
			go s.pool.process(client)
		}
	}
}

// NewServer init new server
func NewServer(config *Config, handler Handler, client ClientConstructor, logger gocli.Logger) *Server {
	opt := options{
		config: config,
		logger: logger,
	}
	return &Server{
		client:  client,
		handler: handler,
		stop:    make(chan struct{}),
		options: opt,
		pool: &pool{
			buffer:      CreateBuffer(int(config.Limits.MaxConnections), int(config.Limits.SharedBufferSize)),
			options:     opt,
			connections: make([]*connection, 0),
		},
	}
}
