package tcpless

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
)

// Signature common interface
type Signature interface {
	// Decode byte message
	Decode(r io.Reader) error
	// Encode byte message
	Encode() []byte
	// Len Length of current message
	Len() uint64
	// Parse current message
	Parse(v any) error
	// RegisterType register custom type
	RegisterType(v any)
	// Send any message
	Send(v any) (response *GobSignature, err error)
	// Stream Get stream
	Stream() io.ReadWriteCloser
}

// GobSignature standard handler signature
type GobSignature struct {
	route  string
	data   []byte
	stream io.ReadWriteCloser
}

// Encode to byte message
func (h *GobSignature) Encode() []byte {
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
	result := make([]byte, uint64(len(h.route))+2+uint64(l2)+ld)
	// copy len of route
	result[0] = l1
	// copy len for len of data
	result[1] = l2
	// copy len of data
	copy(result[2:2+l2], l3[8-l2:])
	// copy route name
	copy(result[2+l2:2+l2+l1], h.route)
	// copy data
	copy(result[2+l2+l1:], h.data)
	// return result
	return result
}

// Decode message from io
func (h *GobSignature) Decode(r io.Reader) error {
	// read route len and len of data len
	l1l2 := [2]byte{}
	_, err := r.Read(l1l2[:])
	if err != nil {
		return err
	}
	// read data len and route
	l := uint64(l1l2[0]) + uint64(l1l2[1])
	l3Route := make([]byte, l)
	_, err = r.Read(l3Route)
	if err != nil {
		return err
	}
	// collect data len
	h.route = string(l3Route[l1l2[1]:])
	l3 := [8]byte{}
	for i := byte(0); i < l1l2[1]; i++ {
		l3[7-i] = l3Route[l1l2[1]-i-1]
	}
	ld := binary.BigEndian.Uint64(l3[:])
	// read data
	h.data = make([]byte, ld)
	_, err = r.Read(h.data)
	return err
}

// Len Length of current message
func (h *GobSignature) Len() uint64 {
	return uint64(len(h.data))
}

// Parse current message
func (h *GobSignature) Parse(v any) error {
	if h.Len() == 0 {
		err := h.Decode(h.stream)
		if err != nil {
			return err
		}
	}
	return gob.NewDecoder(bytes.NewBuffer(h.data)).Decode(v)
}

// RegisterType register custom type message
func (h *GobSignature) RegisterType(v any) {
	gob.Register(v)
}

// Send any message
func (h *GobSignature) Send(v any) (response *GobSignature, err error) {
	b := bytes.NewBuffer(nil)
	err = gob.NewEncoder(b).Encode(v)
	if err != nil {
		return
	}
	h.data = b.Bytes()
	_, err = h.stream.Write(h.Encode())
	if err != nil {
		return
	}
	response = &GobSignature{
		stream: h.stream,
	}
	return
}

// Stream Get stream
func (h *GobSignature) Stream() io.ReadWriteCloser {
	return h.stream
}
