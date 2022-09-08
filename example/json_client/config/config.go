package config

import "github.com/dimonrus/tcpless"

type Config struct {
	TCPLess tcpless.Config `yaml:"tcpLess"`
}
