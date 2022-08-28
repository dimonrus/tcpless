package tcpless

import (
	"errors"
	"io"
	"sync"
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
	return b.data[b.off:]
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

// Each connection have own buffer to encode and decode data
type buffer struct {
	data []*PermanentBuffer
	free []bool
	m    sync.RWMutex
}

// Pull return free buffer
func (b *buffer) Pull() (*PermanentBuffer, uint16) {
	b.m.Lock()
	defer b.m.Unlock()
	for {
		for i, u := range b.free {
			if u {
				b.free[i] = false
				return b.data[i], uint16(i)
			}
		}
	}
}

// Release buffer and clean line
func (b *buffer) Release(index uint16) {
	b.m.Lock()
	defer b.m.Unlock()
	b.data[index].Reset()
	b.free[index] = true
}

// Size of all buffers
func (b *buffer) Size() uint64 {
	b.m.RLock()
	defer b.m.RUnlock()
	var size uint64
	for i := range b.data {
		size += uint64(b.data[i].Cap())
	}
	return size
}

// CreateBuffer create buffer for data
// height - how many connection
// weight - how many bytes required one connection
func CreateBuffer(height, weight int) *buffer {
	b := &buffer{
		data: make([]*PermanentBuffer, height),
		free: make([]bool, height),
		m:    sync.RWMutex{},
	}
	for i := 0; i < height; i++ {
		b.data[i] = NewPermanentBuffer(make([]byte, weight))
		b.free[i] = true
	}
	return b
}
