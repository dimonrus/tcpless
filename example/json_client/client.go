package main

import (
	"bytes"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/porterr"
	"github.com/dimonrus/tcpless"
	"net"
)

// JsonClient transfer json data
type JsonClient struct {
	// Common GetFreeClient
	tcpless.Client
}

// Check for IClient interface
var _ = (tcpless.IClient)(&JsonClient{})

// Dial to server
func (c *JsonClient) Dial() (net.Conn, porterr.IError) {
	conn, err := c.Client.Dial()
	if err != nil {
		return nil, err
	}
	c.SetStream(tcpless.NewConnection(conn, bytes.NewBuffer(make([]byte, 0, tcpless.MinimumSharedBufferSize)), 0))
	return conn, nil
}

// Hello call hello handler
func (c *JsonClient) Hello(user TestUser) (resp TestOkResponse, err error) {
	err = c.Ask("api.hello", user)
	if err != nil {
		return
	}
	err = c.Parse(&resp)
	return
}

// NewJSONClient JSON client constructor
func NewJSONClient(config *tcpless.Config, logger gocli.Logger) tcpless.IClient {
	sig := tcpless.CreateSignature(nil, nil, tcpless.NewJSONDataEncryptor)
	return &JsonClient{Client: tcpless.CreateClient(config, &sig, logger)}
}
