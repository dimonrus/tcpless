package tcpless

import (
	"encoding/binary"
	"io"
)

// Signature common interface
type Signature interface {
	// Data useful message
	Data() []byte
	// Decode byte message
	Decode(r io.Reader, buf *PermanentBuffer) error
	// Encode byte message
	Encode(buf *PermanentBuffer) []byte
	// Len of data
	Len() uint64
	// Route get message route
	Route() string
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
func (h *GobSignature) Decode(r io.Reader, buf *PermanentBuffer) error {
	// reset buffer before nex usage
	defer buf.Reset()
	// read route len and len of data len
	l1l2 := buf.Next(2)
	_, err := r.Read(l1l2)
	if err != nil {
		return err
	}
	// read data len and route
	l := int(l1l2[0]) + int(l1l2[1])
	l3Route := buf.Next(l)
	_, err = r.Read(l3Route)
	if err != nil {
		return err
	}
	// collect data len
	h.route = l3Route[l1l2[1]:]
	l3 := [8]byte{}
	for i := byte(0); i < l1l2[1]; i++ {
		l3[7-i] = l3Route[l1l2[1]-i-1]
	}
	ld := binary.BigEndian.Uint64(l3[:])
	// read data
	h.data = buf.Next(int(ld))
	_, err = r.Read(h.data)
	return err
}

// Encode to byte message
func (h *GobSignature) Encode(buf *PermanentBuffer) []byte {
	defer buf.Reset()
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
	result := buf.Next(len(h.route) + 2 + int(l2) + int(ld))
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

// Len Length of current message
func (h *GobSignature) Len() uint64 {
	return uint64(len(h.data))
}

// Route get route
func (h *GobSignature) Route() string {
	return string(h.route)
}
