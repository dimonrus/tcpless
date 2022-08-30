package tcpless

import (
	"bytes"
	"context"
	"encoding/gob"
	"net"
)

// Check for IClient interface
var _ = (IClient)(&GobClient{})

// ClientConstructor func for specific client init
type ClientConstructor func() IClient

// IClient interface for communication
type IClient interface {
	// Close stream
	Close() error
	// Ctx get context
	Ctx() context.Context
	// Dial to server
	Dial(address net.Addr) error
	// Parse current message
	Parse(v any) error
	// RegisterType register custom type
	RegisterType(v any)
	// Send any message
	Send(route string, v any) error
	// Read get signature from stream
	Read() (Signature, error)
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
}

// Close stream
func (c *Client) Close() error {
	return c.stream.Release()
}

// Ctx get context
func (c *Client) Ctx() context.Context {
	return c.ctx
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

// Dial to server
func (g *GobClient) Dial(address net.Addr) error {
	conn, err := net.Dial(address.Network(), address.String())
	if err != nil {
		return err
	}
	// TODO permanent buffer size
	g.SetStream(&connection{
		Conn:   conn,
		done:   make(chan struct{}),
		buffer: bytes.NewBuffer(make([]byte, 0, 1024)),
		index:  0,
	})
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
	return
}

// Send data to stream
func (g *GobClient) Send(route string, v any) error {
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

// NewGobClient gob client constructor
func NewGobClient() IClient {
	return &GobClient{Client: Client{sig: &GobSignature{}}}
}
