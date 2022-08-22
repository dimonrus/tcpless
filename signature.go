package tcpless

import "io"

// Signature common interface
type Signature interface {
	Auth() bool
	Route() string
	Len() uint64
	Stream() io.ReadWriteCloser
}

// HandlerData standard handler signature
type HandlerData struct {
	route         string
	contentLength uint64
	stream        io.ReadWriteCloser
}

func (h HandlerData) Auth() bool {
	return true
}

func (h HandlerData) Len() uint64 {
	return h.contentLength
}

func (h HandlerData) Route() string {
	return h.route
}

func (h HandlerData) Stream() io.ReadWriteCloser {
	return h.stream
}
