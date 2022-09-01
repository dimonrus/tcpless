package tcpless

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/dimonrus/gocli"
	"net"
)

// Check for IClient interface
var _ = (IClient)(&Client{})

// ClientConstructor func for specific client init
type ClientConstructor func(config *Config, logger gocli.Logger) IClient

// IClient interface for communication
type IClient interface {
	// Ask any message
	Ask(route string, v any) error
	// Close stream
	Close() error
	// Ctx get context
	Ctx() context.Context
	// Dial to server
	Dial() (net.Conn, error)
	// Logger get logger
	Logger() gocli.Logger
	// Parse current message
	Parse(v any) error
	// Read get signature from stream
	Read() (Signature, error)
	// RegisterType register custom type
	RegisterType(v any)
	// SetStream set stream io
	SetStream(stream Connection)
	// Stream Get stream
	Stream() Connection
	// WithContext With context
	WithContext(ctx context.Context)
}

// Client structure
type Client struct {
	// connection
	stream Connection
	// signature
	sig Signature
	// context
	ctx context.Context
	// options
	options options
}

// Ask any message
func (c *Client) Ask(route string, v any) error {
	// Nothing to implements
	return errors.New("no action in base Client")
}

// Dial to server
func (c *Client) Dial() (net.Conn, error) {
	if c.options.config.TLS.Enabled {
		config, err := c.options.config.TLS.LoadTLSConfig()
		if err != nil {
			return nil, err
		}
		return tls.Dial(c.options.config.Address.Network(), c.options.config.Address.String(), config)
	}
	return net.Dial(c.options.config.Address.Network(), c.options.config.Address.String())
}

// Close stream
func (c *Client) Close() error {
	return c.stream.Release()
}

// Ctx get context
func (c *Client) Ctx() context.Context {
	return c.ctx
}

// Logger get logger
func (c *Client) Logger() gocli.Logger {
	return c.options.logger
}

// Parse current message into v
func (c *Client) Parse(v any) error {
	// Nothing to implements
	return errors.New("no action in Client for method Parse")
}

// Read get signature from stream
func (c *Client) Read() (Signature, error) {
	// Nothing to implements
	return c.sig, errors.New("no action in Client for method Read")
}

// RegisterType register custom type
func (c *Client) RegisterType(v any) {
	// Nothing to implements
}

// SetStream set stream io
func (c *Client) SetStream(stream Connection) {
	// set stream
	c.stream = stream
}

// Stream Get stream
func (c *Client) Stream() Connection {
	return c.stream
}

// WithContext with context
func (c *Client) WithContext(ctx context.Context) {
	c.ctx = ctx
	return
}
