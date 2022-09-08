package tcpless

import (
	"bytes"
	"github.com/dimonrus/gocli"
	"net"
)

// Check for IClient interface
var _ = (IClient)(&GobClient{})

// GobClient GetFreeClient for gob serialization
type GobClient struct {
	// Common GetFreeClient
	Client
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

// NewGobClient gob GetFreeClient constructor
func NewGobClient(config *Config, logger gocli.Logger) IClient {
	// Create signature with gob encryption/decryption
	sig := CreateSignature(nil, nil, NewGobDataEncryptor)
	return &GobClient{Client: CreateClient(config, &sig, logger)}
}
