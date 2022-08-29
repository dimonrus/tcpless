package tcpless

import (
	"bytes"
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
}

// Client structure
type Client struct {
	// connection
	stream Connection
	// signature
	sig Signature
}

// Close stream
func (c *Client) Close() error {
	return c.stream.Release()
}

// Dial to server
func (c *Client) Dial(address net.Addr) error {
	conn, err := net.Dial(address.Network(), address.String())
	if err != nil {
		return err
	}
	// TODO permanent buffer size
	c.stream = &connection{
		Conn:   conn,
		done:   make(chan struct{}),
		buffer: bytes.NewBuffer(make([]byte, 0, 1024)),
		index:  0,
	}
	return err
}

// Stream Get stream
func (c *Client) Stream() Connection {
	return c.stream
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

// Parse data to type
func (g *GobClient) Parse(v any) error {
	var err error
	if g.sig.Len() == 0 {
		err = g.sig.Decode(g.stream.Connection(), g.stream.Buffer())
		if err != nil {
			return err
		}
	}
	return g.decoder.Decode(v)
}

// RegisterType register type for communication
func (g *GobClient) RegisterType(v any) {
	gob.Register(v)
	return
}

// Signature get from stream
func (g *GobClient) Read() (Signature, error) {
	err := g.sig.Decode(g.stream.Connection(), g.stream.Buffer())
	return g.sig, err
}

// SetStream set stream
func (g *GobClient) SetStream(stream Connection) {
	// set stream
	g.stream = stream
	// set encoder
	g.encoder = gob.NewEncoder(stream.Buffer())
	// set stream
	g.decoder = gob.NewDecoder(stream.Buffer())
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
	_, err = g.stream.Connection().Write(s.Encode(g.stream.Buffer()))
	return err
}

// NewGobClient gob client constructor
func NewGobClient() IClient {
	return &GobClient{Client: Client{sig: &GobSignature{}}}
}
