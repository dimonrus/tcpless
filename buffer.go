package tcpless

import (
	"bytes"
	"sync"
)

type Buffer struct {
	data []byte
	off  uint64
}

func (b *Buffer) Next(n uint64) []byte {
	res := b.data[b.off : b.off+n]
	b.off += n
	return res
}

func (b *Buffer) Reset() {
	b.off = 0
	return
}

func (b *Buffer) Grow(n uint64) {
	if len(b.data) == cap(b.data) {

	}
	return
}

// Each connection have own buffer to encode and decode data
type buffer struct {
	data []*bytes.Buffer
	free []bool
	m    sync.RWMutex
}

// Pull return free buffer
func (b *buffer) Pull() (*bytes.Buffer, uint16) {
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
	b.data[index].Bytes()
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
		data: make([]*bytes.Buffer, height),
		free: make([]bool, height),
		m:    sync.RWMutex{},
	}
	for i := 0; i < height; i++ {
		b.data[i] = bytes.NewBuffer(make([]byte, weight))
		b.free[i] = true
	}
	return b
}
