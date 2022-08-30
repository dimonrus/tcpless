package tcpless

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/dimonrus/gocli"
	"net"
)

// Check for IClient interface
var _ = (IClient)(&GobClient{})

// ClientConstructor func for specific client init
type ClientConstructor func(logger gocli.Logger) IClient

// IClient interface for communication
type IClient interface {
	// Ask any message
	Ask(route string, v any) error
	// Close stream
	Close() error
	// Ctx get context
	Ctx() context.Context
	// Dial to server
	Dial(address net.Addr) error
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
	// client logger
	logger gocli.Logger
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
	return c.logger
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

// GobClient client for gob serialization
type GobClient struct {
	// Common client
	Client
	// Gob decoder
	decoder *gob.Decoder
	// Gob encoder
	encoder *gob.Encoder
}

// Ask server
func (g *GobClient) Ask(route string, v any) error {
	g.stream.Buffer().Reset()
	err := g.encoder.Encode(v)
	if err != nil {
		return err
	}
	s := GobSignature{
		route: []byte(route),
		data:  g.stream.Buffer().Bytes(),
	}
	g.stream.Buffer().Reset()
	_, err = g.stream.Connection().Write(s.Encode(g.stream.Buffer()))
	return err
}

// Dial to server
func (g *GobClient) Dial(address net.Addr) error {
	conn, err := net.Dial(address.Network(), address.String())
	if err != nil {
		return err
	}
	g.SetStream(newConnection(conn, bytes.NewBuffer(make([]byte, 0, MinimumSharedBufferSize)), 0))
	return err
}

// Parse data to type
func (g *GobClient) Parse(v any) error {
	var err error
	// if signature is empty
	if g.sig.Len() == 0 {
		// read data for io
		_, err = g.Read()
		if err != nil {
			return err
		}
	}
	// clear buffer after read
	g.stream.Buffer().Reset()
	// write in buffer only data
	g.stream.Buffer().Write(g.sig.Data())
	// Reset signature
	g.sig.Reset()
	// decode value
	return g.decoder.Decode(v)
}

// Signature get from stream
func (g *GobClient) Read() (Signature, error) {
	// clear buffer before
	g.stream.Buffer().Reset()
	// decode message to signature
	err := g.sig.Decode(g.stream.Connection(), g.stream.Buffer())
	return g.sig, err
}

// RegisterType register type for communication
func (g *GobClient) RegisterType(v any) {
	gob.Register(v)
	return
}

// SetStream set stream
func (g *GobClient) SetStream(stream Connection) {
	// set stream
	g.stream = stream
	// set stream
	g.decoder = gob.NewDecoder(stream.Buffer())
	// set encoder
	g.encoder = gob.NewEncoder(stream.Buffer())
}

// NewGobClient gob client constructor
func NewGobClient(logger gocli.Logger) IClient {
	return &GobClient{Client: Client{sig: &GobSignature{}, logger: logger}}
}
