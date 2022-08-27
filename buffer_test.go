package tcpless

import (
	"testing"
)

func TestCreateBuffer(t *testing.T) {
	buf := CreateBuffer(20, 100)
	if buf.Size() != 2000 {
		t.Fatal("wrong create buffer")
	}
}

func TestBuffer_Pull(t *testing.T) {
	buf := CreateBuffer(20, 100)
	b, i := buf.Pull()
	b.Reset()
	b.Write([]byte("HelloWorld"))
	buf.Release(i)
	if buf.Size() != 2000 {
		t.Fatal("wrong size")
	}
}

func BenchmarkBuffer_Pull(b *testing.B) {
	buf := CreateBuffer(20, 100)
	bt, _ := buf.Pull()
	bt.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Write([]byte("HelloWorld"))
		bt.Reset()
	}
	b.ReportAllocs()
}
