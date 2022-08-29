package tcpless

import (
	"bytes"
	"sync"
)

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
		b.data[i] = bytes.NewBuffer(make([]byte, 0, weight))
		b.free[i] = true
	}
	return b
}
