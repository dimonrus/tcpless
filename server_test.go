package tcpless

import (
	"net"
	"time"
)

func getTestConfig() *Config {
	return &Config{
		Address: &net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 900,
		},
		Limits: ConnectionLimit{
			MaxConnections:   10,
			SharedBufferSize: 1024,
			MaxIdle:          time.Second * 10,
			RedialTimeout:    time.Second,
		},
	}
}
