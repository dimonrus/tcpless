package tcpless

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"sync"
)

// TLSConfig configuration
type TLSConfig struct {
	// Is enabled
	Enabled bool `yaml:"enabled"`
	// Root CA path
	CaPath string `yaml:"caPath"`
	// Path to cert
	CertPath string `yaml:"certPath"`
	// Path to key
	KeyPath string `yaml:"keyPath"`
	// lazy load
	config *tls.Config
	// once load config
	m sync.Once
}

// LoadTLSConfig load and prepare tls.Config
func (c *TLSConfig) LoadTLSConfig() (*tls.Config, error) {
	if c.config != nil {
		return c.config, nil
	}
	var err error
	c.m.Do(func() {
		var cert tls.Certificate
		cert, err = tls.LoadX509KeyPair(c.CertPath, c.KeyPath)
		if err != nil {
			return
		}
		c.config = &tls.Config{Certificates: []tls.Certificate{cert}}
		if c.CaPath == "" {
			return
		}
		var caCert []byte
		caCert, err = os.ReadFile(c.CaPath)
		if err != nil {
			return
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		c.config.RootCAs = caCertPool
	})
	return c.config, nil
}
