package tcpless

import (
	"errors"
	"io"
)

// PermanentBuffer local buffer
type PermanentBuffer struct {
	// data
	data []byte
	// offset
	off int
}

// Next get n bytes
func (b *PermanentBuffer) Next(n int) []byte {
	if n+b.off > cap(b.data) {
		panic("n + offset > cap(data); check your shared buffer size")
		return nil
	}
	res := b.data[b.off : b.off+n]
	b.off += n
	return res
}

// Seek change offset
func (b *PermanentBuffer) Seek(offset int64, whence int) (int64, error) {
	var abs int
	switch whence {
	case io.SeekStart:
		abs = int(offset)
	case io.SeekCurrent:
		abs = b.off + int(offset)
	case io.SeekEnd:
		abs = len(b.data) + int(offset)
	default:
		return 0, errors.New("tcpless.PermanentBuffer.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("tcpless.PermanentBuffer.Seek: negative position")
	}
	b.off = abs
	return int64(abs), nil
}

// Reset cursor
func (b *PermanentBuffer) Reset() {
	b.off = 0
	return
}

// Cap return cap for buffer
func (b *PermanentBuffer) Cap() int {
	return cap(b.data)
}

// Bytes return data
func (b *PermanentBuffer) Bytes() []byte {
	return b.data[:b.off]
}

// Read bytes
func (b *PermanentBuffer) Read(p []byte) (n int, err error) {
	if b.off >= b.Cap() {
		return 0, io.EOF
	}
	n = copy(p, b.data[b.off:])
	b.off += n
	return n, nil
}

// Bytes return data
func (b *PermanentBuffer) Write(data []byte) (n int, err error) {
	err = nil
	if len(data)+b.off > cap(b.data) {
		panic("len(data) + offset > cap(data)")
		return
	}
	n = copy(b.data[b.off:], data[:])
	b.off += n
	return
}

// NewPermanentBuffer init buffer
func NewPermanentBuffer(data []byte) *PermanentBuffer {
	return &PermanentBuffer{
		data: data,
		off:  0,
	}
}
