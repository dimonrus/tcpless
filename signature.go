package tcpless

import (
	"bytes"
	"encoding/binary"
	"io"
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
	// reset buffer before next usage
	// defer buf.Reset()
	// read all from reader
	data := buf.Bytes()
	n, err := r.Read(data[:buf.Cap()])
	if err != nil {
		return err
	}
	data = data[:n]
	// read route len and len of data len
	l1l2 := data[:2]
	// read data len and route
	l := int(l1l2[0]) + int(l1l2[1])
	l3Route := data[2 : 2+l]
	// collect data len
	h.route = l3Route[l1l2[1]:]
	l3 := [8]byte{}
	for i := byte(0); i < l1l2[1]; i++ {
		l3[7-i] = l3Route[l1l2[1]-i-1]
	}
	ld := binary.BigEndian.Uint64(l3[:])
	// read data
	h.data = data[2+l : 2+l+(int(ld))]
	// write data to buf
	buf.Write(data)
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
	var l1 = byte(len(h.route))
	// len of byte for data len
	var l2 byte
	// data length
	l3 := [8]byte{}
	ld := uint64(len(h.data))
	for i := 7; i >= 0; i-- {
		if ld>>(i*8) > 0 {
			l2++
			l3[7-i] = byte(ld >> (i * 8))
		}
	}
	// make result slice
	result := buf.Bytes()[:(len(h.route) + 2 + int(l2) + int(ld))]
	// copy data. Do it before route name will be saved
	copy(result[2+l2+l1:], h.data)
	// copy len of route
	result[0] = l1
	// copy len for len of data
	result[1] = l2
	// copy len of data
	copy(result[2:2+l2], l3[8-l2:])
	// copy route name
	copy(result[2+l2:2+l2+l1], h.route)
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
	return string(h.route)
}

// Write rewrite bytes
func (h *GobSignature) Write(p []byte) (n int, err error) {
	n = copy(h.data[:], p)
	return
}
