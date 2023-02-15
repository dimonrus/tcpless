package tcpless

import (
	"bytes"
	"context"
	"crypto/tls"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/porterr"
	"net"
	"time"
)

// Check for IClient interface
var _ = (IClient)(&Client{})

// ClientConstructor func for specific GetFreeClient init
type ClientConstructor func(config *Config, logger gocli.Logger) IClient

// IClient interface for communication
type IClient interface {
	// Ask custom type message
	Ask(route string, v any) porterr.IError
	// AskBytes send bytes
	AskBytes(route string, b []byte) porterr.IError
	// Ctx get context
	Ctx() context.Context
	// Dial to server
	Dial() (net.Conn, porterr.IError)
	// Logger get logger
	Logger() gocli.Logger
	// Parse current message
	Parse(v any) porterr.IError
	// Read get signature from stream
	Read() (ISignature, porterr.IError)
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
func (c *Client) Ask(route string, v any) porterr.IError {
	// reset buffer
	c.stream.Buffer().Reset()
	if e, ok := v.(porterr.IError); ok {
		e.Origin().Pack(c.stream.Buffer())
		return c.AskBytes(route, c.stream.Buffer().Bytes())
	}
	for {
		// encode v to bytes
		err := c.sig.Encryptor().Encode(v)
		if err != nil {
			return porterr.New(porterr.PortErrorEncoder, err.Error())
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
		if err == nil {
			return nil
		} else {
			e := c.reDial(err)
			if e != nil {
				return e
			}
		}
	}
}

// redial
func (c *Client) reDial(err error) porterr.IError {
	if oe, ok := err.(*net.OpError); ok && (oe.Op == "write") {
		for {
			time.Sleep(c.options.config.Limits.RedialTimeout * time.Second)
			conn, e := c.Dial()
			if e != nil && c.Logger() != nil {
				c.Logger().Errorln(e.Error())
			} else {
				c.SetStream(NewConnection(conn, bytes.NewBuffer(make([]byte, 0, MinimumSharedBufferSize)), 0))
				return nil
			}
		}
	}
	return porterr.New(porterr.PortErrorIO, err.Error())
}

// AskBytes send bytes
func (c *Client) AskBytes(route string, b []byte) porterr.IError {
	// reset buffer
	c.stream.Buffer().Reset()
	// create signature
	s := Signature{
		route: []byte(route),
		data:  b,
	}
	for {
		_, err := c.stream.Connection().Write(s.Encode(c.stream.Buffer()))
		if err == nil {
			return nil
		} else {
			e := c.reDial(err)
			if e != nil {
				return e
			}
		}
	}
}

// Dial to server
func (c *Client) Dial() (net.Conn, porterr.IError) {
	var err error
	var conn net.Conn
	if c.options.config.TLS.Enabled {
		var config *tls.Config
		config, err = c.options.config.TLS.LoadTLSConfig()
		if err != nil {
			return nil, porterr.New(porterr.PortErrorConflict, err.Error())
		}
		conn, err = tls.Dial(c.options.config.Address.Network(), c.options.config.Address.String(), config)
		if err != nil {
			return nil, porterr.New(porterr.PortErrorConnection, err.Error())
		}
		return conn, nil
	}
	conn, err = net.Dial(c.options.config.Address.Network(), c.options.config.Address.String())
	if err != nil {
		return nil, porterr.New(porterr.PortErrorConnection, err.Error())
	}
	return conn, nil
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
func (c *Client) Parse(v any) porterr.IError {
	var e porterr.IError
	// if signature is empty
	if c.sig.Len() == 0 {
		// read data for io
		_, err := c.Read()
		if err != nil {
			return porterr.New(porterr.PortErrorIO, err.Error())
		}
	}
	// clear buffer after read
	c.stream.Buffer().Reset()
	// check if portable error
	e = porterr.UnPack(c.sig.Data())
	if e != nil {
		return e
	}
	// write in buffer only data
	c.stream.Buffer().Write(c.sig.Data())
	// Reset signature
	c.sig.Reset()
	// decode value
	err := c.Signature().Encryptor().Decode(v)
	if err != nil {
		return porterr.New(porterr.PortErrorDecoder, err.Error())
	}
	return nil
}

// Read get signature from stream
func (c *Client) Read() (ISignature, porterr.IError) {
	// clear buffer before
	c.stream.Buffer().Reset()
	// decode message to signature
	err := c.sig.Decode(c.stream.Connection(), c.stream.Buffer())
	if err != nil {
		return nil, porterr.New(porterr.PortErrorDecoder, err.Error())
	}
	// return signature and error
	return c.sig, nil
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
