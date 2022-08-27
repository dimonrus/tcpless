package tcpless

import (
	"bytes"
	"encoding/gob"
	"net"
)

// Check for IClient interface
var _ = (IClient)(&GobClient{})

// IClient interface for communication
type IClient interface {
	// Close stream
	Close() error
	// Dial to server
	Dial(address net.Addr) error
	// Parse current message
	Parse(signature Signature, v any) error
	// RegisterType register custom type
	RegisterType(v any)
	// Send any message
	Send(route string, v any) error
	// Read get signature from stream
	Read() Signature
	// SetStream set stream io
	SetStream(stream Connection)
	// Stream Get stream
	Stream() Connection
}

// Client structure
type Client struct {
	// connection
	stream Connection
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
	c.stream = &connection{
		Conn:   conn,
		done:   make(chan struct{}),
		buffer: nil,
		index:  0,
	}
	return err
}

// SetStream set stream
func (c *Client) SetStream(stream Connection) {
	c.stream = stream
	return
}

// Stream Get stream
func (c *Client) Stream() Connection {
	return c.stream
}

// GobClient client for gob serialization
type GobClient struct {
	Client
}

// Parse data to type
func (g *GobClient) Parse(signature Signature, v any) error {
	var err error
	if signature.Len() == 0 {
		err = signature.Decode(g.stream.Connection(), g.stream.Buffer())
		if err != nil {
			return err
		}
	}
	return gob.NewDecoder(bytes.NewBuffer(signature.Data())).Decode(v)
}

// RegisterType register type for communication
func (g *GobClient) RegisterType(v any) {
	gob.Register(v)
	return
}

// Signature get from stream
func (g *GobClient) Read() Signature {
	sig := &GobSignature{}
	err := sig.Decode(g.stream.Connection(), g.stream.Buffer())
	if err != nil {
		return nil
	}
	return sig
}

// Send data to stream
func (g *GobClient) Send(route string, v any) error {
	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(v)
	if err != nil {
		return err
	}
	s := GobSignature{
		route: []byte(route),
		data:  b.Bytes(),
	}
	_, err = g.stream.Connection().Write(s.Encode(g.stream.Buffer()))
	return err
}
