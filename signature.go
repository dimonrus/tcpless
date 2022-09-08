package tcpless

import (
	"bytes"
	"encoding/binary"
	"io"
	"unsafe"
)

// Check Signature for ISignature interface
var _ = (ISignature)(&Signature{})

// ISignature common interface
type ISignature interface {
	// Data useful message
	Data() []byte
	// Decode byte message
	Decode(r io.Reader, buf *bytes.Buffer) error
	// Encode byte message
	Encode(buf *bytes.Buffer) []byte
	// Encryptor get data encryptor
	Encryptor() DataEncryptor
	// InitEncryptor return init data encryptor constructor
	InitEncryptor(buf *bytes.Buffer)
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

// Signature standard handler signature
type Signature struct {
	// route
	route []byte
	// data
	data []byte
	// apply with buffer on set stream method
	initEncryptor DataEncryptorConstructor
	// data encryptor
	encryptor DataEncryptor
}

// Data get useful message
func (h *Signature) Data() []byte {
	return h.data
}

// Decode message from io
// r - input with bytes
// buf - bytes buffer
func (h *Signature) Decode(r io.Reader, buf *bytes.Buffer) error {
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
func (h *Signature) Encode(buf *bytes.Buffer) []byte {
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

// Encryptor get data encryptor
func (h *Signature) Encryptor() DataEncryptor {
	return h.encryptor
}

// InitEncryptor init encryptor
func (h *Signature) InitEncryptor(buf *bytes.Buffer) {
	h.encryptor = h.initEncryptor(buf)
	return
}

// Len Length of current message
func (h *Signature) Len() uint64 {
	return uint64(len(h.data))
}

// Read all bytes
func (h *Signature) Read(p []byte) (n int, err error) {
	n = copy(p, h.data[:])
	return
}

// Reset signature
func (h *Signature) Reset() {
	h.route = nil
	h.data = nil
	return
}

// Route get route
func (h *Signature) Route() string {
	return *(*string)(unsafe.Pointer(&h.route))
}

// Write rewrite bytes
func (h *Signature) Write(p []byte) (n int, err error) {
	n = copy(h.data[:], p)
	return
}

// CreateSignature prepare Signature struct
func CreateSignature(route []byte, data []byte, initEncryptor DataEncryptorConstructor) Signature {
	return Signature{
		route:         route,
		data:          data,
		initEncryptor: initEncryptor,
	}
}
