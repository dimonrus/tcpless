package tcpless

import (
	"context"
	"crypto/tls"
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

// Ask server route
func (c *Client) Ask(route string, v any) error {
	// reset buffer
	c.stream.Buffer().Reset()
	// encode v to bytes
	err := c.sig.Encryptor().Encode(v)
	if err != nil {
		return err
	}
	// create signature
	s := Signature{
		route: []byte(route),
		data:  c.stream.Buffer().Bytes(),
	}
	// reset buffer
	c.stream.Buffer().Reset()
	// send data to stream
	_, err = c.stream.Connection().Write(s.Encode(c.stream.Buffer()))
	return err
}

// AskBytes send bytes
func (c *Client) AskBytes(route string, b []byte) error {
	// reset buffer
	c.stream.Buffer().Reset()
	// create signature
	s := Signature{
		route: []byte(route),
		data:  b,
	}
	// write bytes into connection
	_, err := c.stream.Connection().Write(s.Encode(c.stream.Buffer()))
	return err
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

// Ctx get context
func (c *Client) Ctx() context.Context {
	return c.ctx
}

// Logger get logger
func (c *Client) Logger() gocli.Logger {
	return c.options.logger
}

// Parse data to type
func (c *Client) Parse(v any) error {
	var err error
	// if signature is empty
	if c.sig.Len() == 0 {
		// read data for io
		_, err = c.Read()
		if err != nil {
			return err
		}
	}
	// clear buffer after read
	c.stream.Buffer().Reset()
	// write in buffer only data
	c.stream.Buffer().Write(c.sig.Data())
	// Reset signature
	c.sig.Reset()
	// decode value
	return c.Signature().Encryptor().Decode(v)
}

// Read get signature from stream
func (c *Client) Read() (ISignature, error) {
	// clear buffer before
	c.stream.Buffer().Reset()
	// decode message to signature
	err := c.sig.Decode(c.stream.Connection(), c.stream.Buffer())
	// return signature and error
	return c.sig, err
}

// SetStream set stream io
func (c *Client) SetStream(stream Streamer) {
	// set stream
	c.stream = stream
	// set up encryptor
	c.sig.InitEncryptor(stream.Buffer())
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
