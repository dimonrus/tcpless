package tcpless

import (
	"bytes"
	"encoding/gob"
	"github.com/dimonrus/gocli"
	"net"
)

// Check for IClient interface
var _ = (IClient)(&GobClient{})

// GobClient GetFreeClient for gob serialization
type GobClient struct {
	// Common GetFreeClient
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
	s := Signature{
		route: []byte(route),
		data:  g.stream.Buffer().Bytes(),
	}
	g.stream.Buffer().Reset()
	_, err = g.stream.Connection().Write(s.Encode(g.stream.Buffer()))
	return err
}

// AskBytes send bytes
func (g *GobClient) AskBytes(route string, b []byte) error {
	g.stream.Buffer().Reset()
	s := Signature{
		route: []byte(route),
		data:  b,
	}
	g.stream.Buffer().Reset()
	_, err := g.stream.Connection().Write(s.Encode(g.stream.Buffer()))
	return err
}

// Dial to server
func (g *GobClient) Dial() (net.Conn, error) {
	conn, err := g.Client.Dial()
	if err != nil {
		return nil, err
	}
	g.SetStream(NewConnection(conn, bytes.NewBuffer(make([]byte, 0, MinimumSharedBufferSize)), 0))
	return conn, nil
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

// ISignature get from stream
func (g *GobClient) Read() (ISignature, error) {
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
func (g *GobClient) SetStream(stream Streamer) {
	// set stream
	g.stream = stream
	// set stream
	g.decoder = gob.NewDecoder(stream.Buffer())
	// set encoder
	g.encoder = gob.NewEncoder(stream.Buffer())
}

// NewGobClient gob GetFreeClient constructor
func NewGobClient(config *Config, logger gocli.Logger) IClient {
	return &GobClient{Client: CreateClient(config, &Signature{}, logger)}
}
