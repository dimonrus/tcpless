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

// ClientConstructor func for specific GetFreeClient init
type ClientConstructor func(config *Config, logger gocli.Logger) IClient

// IClient interface for communication
type IClient interface {
	// Ask custom type message
	Ask(route string, v any) error
	// AskBytes send bytes
	AskBytes(route string, b []byte) error
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
	Read() (ISignature, error)
	// RegisterType register custom type
	RegisterType(v any)
	// SetStream set stream io
	SetStream(stream Streamer)
	// Signature return signature
	Signature() ISignature
	// Stream Get stream
	Stream() Streamer
	// WithContext With context
	WithContext(ctx context.Context)
}

// Client structure
type Client struct {
	// connection
	stream Streamer
	// signature
	sig ISignature
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

// AskBytes any message
func (c *Client) AskBytes(route string, b []byte) error {
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
func (c *Client) Read() (ISignature, error) {
	// Nothing to implements
	return c.sig, errors.New("no action in Client for method Read")
}

// RegisterType register custom type
func (c *Client) RegisterType(v any) {
	// Nothing to implements
}

// SetStream set stream io
func (c *Client) SetStream(stream Streamer) {
	// set stream
	c.stream = stream
}

// Signature return signature
func (c *Client) Signature() ISignature {
	return c.sig
}

// Stream Get stream
func (c *Client) Stream() Streamer {
	return c.stream
}

// WithContext with context
func (c *Client) WithContext(ctx context.Context) {
	c.ctx = ctx
	return
}

// CreateClient create base client
func CreateClient(config *Config, sig ISignature, logger gocli.Logger) Client {
	return Client{sig: sig, options: options{config: config, logger: logger}}
}
