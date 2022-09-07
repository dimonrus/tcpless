package json_client

import (
	"bytes"
	"encoding/json"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/tcpless"
	"net"
)

// Check for IClient interface
var _ = (tcpless.IClient)(&JsonClient{})

// JsonClient transfer json data
type JsonClient struct {
	// Common GetFreeClient
	tcpless.Client
	// Gob decoder
	decoder *json.Decoder
	// Gob encoder
	encoder *json.Encoder
}

// Ask server
func (c *JsonClient) Ask(route string, v any) error {
	// reset buffer
	c.Stream().Buffer().Reset()
	// json encode
	err := c.encoder.Encode(v)
	if err != nil {
		return err
	}
	// create signature
	s := tcpless.CreateSignature([]byte(route), c.Stream().Buffer().Bytes())
	// reset buffer
	c.Stream().Buffer().Reset()
	// send data to stream
	_, err = c.Stream().Connection().Write(s.Encode(c.Stream().Buffer()))
	return err
}

// AskBytes send bytes
func (c *JsonClient) AskBytes(route string, b []byte) error {
	// reset buffer
	c.Stream().Buffer().Reset()
	// create signature
	s := tcpless.CreateSignature([]byte(route), b)
	// send data to stream
	_, err := c.Stream().Connection().Write(s.Encode(c.Stream().Buffer()))
	return err
}

// Dial to server
func (c *JsonClient) Dial() (net.Conn, error) {
	conn, err := c.Client.Dial()
	if err != nil {
		return nil, err
	}
	c.SetStream(tcpless.NewConnection(conn, bytes.NewBuffer(make([]byte, 0, tcpless.MinimumSharedBufferSize)), 0))
	return conn, nil
}

// Parse data to type
func (c *JsonClient) Parse(v any) error {
	var err error
	// if signature is empty
	if c.Signature().Len() == 0 {
		// read data for io
		_, err = c.Read()
		if err != nil {
			return err
		}
	}
	// clear buffer after read
	c.Stream().Buffer().Reset()
	// write in buffer only data
	c.Stream().Buffer().Write(c.Signature().Data())
	// Reset signature
	c.Signature().Reset()
	// decode value
	return c.decoder.Decode(v)
}

// ISignature get from stream
func (c *JsonClient) Read() (tcpless.ISignature, error) {
	// clear buffer before
	c.Stream().Buffer().Reset()
	// decode message to signature
	err := c.Signature().Decode(c.Stream().Connection(), c.Stream().Buffer())
	// return signature and error if exists
	return c.Signature(), err
}

// SetStream set stream
func (c *JsonClient) SetStream(stream tcpless.Streamer) {
	// set stream
	c.Client.SetStream(stream)
	// set stream
	c.decoder = json.NewDecoder(stream.Buffer())
	// set encoder
	c.encoder = json.NewEncoder(stream.Buffer())
}

// NewJSONClient JSON client constructor
func NewJSONClient(config *tcpless.Config, logger gocli.Logger) tcpless.IClient {
	return &JsonClient{
		Client: tcpless.CreateClient(config, &tcpless.Signature{}, logger),
	}
}
