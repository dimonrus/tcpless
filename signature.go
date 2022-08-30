package tcpless

import (
	"bytes"
	"encoding/binary"
	"io"
	"unsafe"
)

// Signature common interface
type Signature interface {
	// Data useful message
	Data() []byte
	// Decode byte message
	Decode(r io.Reader, buf *bytes.Buffer) error
	// Encode byte message
	Encode(buf *bytes.Buffer) []byte
	// Len of data
	Len() uint64
	// Read implements reader interface
	Read(p []byte) (n int, err error)
	// Reset signature
	Reset()
	// Route get message route
	Route() string
	// Write implements writer interface
	Write(p []byte) (n int, err error)
}

// GobSignature standard handler signature
type GobSignature struct {
	// route
	route []byte
	// data
	data []byte
}

// Data get useful message
func (h *GobSignature) Data() []byte {
	return h.data
}

// Decode message from io
// r - input with bytes
// buf - bytes buffer
func (h *GobSignature) Decode(r io.Reader, buf *bytes.Buffer) error {
	// store data
	data := buf.Bytes()
	// read first 4 bytes
	_, err := r.Read(data[:4])
	if err != nil {
		return err
	}
	// header length
	l := data[:4]
	// len of route
	l1 := l[0]
	// set l[0] to zero
	l[0] = 0
	// get len of data
	l2 := binary.BigEndian.Uint32(l[:])
	// must read content
	n := int(l1) + int(l2)
	// read content
	n, err = r.Read(data[:n])
	if err != nil {
		return err
	}
	// read data
	h.data = data[l1:n]
	// set route
	h.route = data[:l1]
	return nil
}

// Encode to byte message
func (h *GobSignature) Encode(buf *bytes.Buffer) []byte {
	// reset buffer before next usage
	buf.Reset()
	// route length
	if len(h.route) > 255 {
		return nil
	}
	// header of encoded package
	// l[0] - len of route
	// l[1:] - len of data
	l := [4]byte{}
	// len of route
	l[0] = byte(len(h.route))
	// create len of data slice
	ld := uint32(len(h.data))
	for i := 3; i >= 0; i-- {
		if ld>>(i*8) > 0 {
			l[3-i] = byte(ld >> (i * 8))
		}
	}
	// make result slice
	result := buf.Bytes()[:(len(h.route) + 4 + int(ld))]
	// copy data. Do it before route name will be saved
	copy(result[4+int(l[0]):], h.data)
	// copy route name
	copy(result[4:4+int(l[0])], h.route)
	// copy header
	copy(result[:4], l[:])
	// write to buffer
	buf.Write(result)
	return buf.Bytes()
}

// Len Length of current message
func (h *GobSignature) Len() uint64 {
	return uint64(len(h.data))
}

// Read all bytes
func (h *GobSignature) Read(p []byte) (n int, err error) {
	n = copy(p, h.data[:])
	return
}

// Reset signature
func (h *GobSignature) Reset() {
	h.route = nil
	h.data = nil
	return
}

// Route get route
func (h *GobSignature) Route() string {
	return *(*string)(unsafe.Pointer(&h.route))
}

// Write rewrite bytes
func (h *GobSignature) Write(p []byte) (n int, err error) {
	n = copy(h.data[:], p)
	return
}
