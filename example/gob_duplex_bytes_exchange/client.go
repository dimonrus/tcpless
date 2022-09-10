package main

import (
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/tcpless"
)

// HelloClient transfer пщи data
type HelloClient struct {
	// Common GetFreeClient
	tcpless.IClient
}

// Hello call hello handler
func (c *HelloClient) Hello(data []byte) error {
	return c.AskBytes("api.v1.exchange", data)
}

// NewHelloClient client constructor
func NewHelloClient(config *tcpless.Config, logger gocli.Logger) tcpless.IClient {
	return &HelloClient{IClient: tcpless.NewGobClient(config, logger)}
}
