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
	n, _ := b.Write([]byte("HelloWorld"))
	if n != 10 {
		t.Fatal("wrong write")
	}
	buf.Release(i)
	if buf.Size() != 2000 {
		t.Fatal("wrong size")
	}
}

func BenchmarkBuffer_Pull(b *testing.B) {
	buf := CreateBuffer(20, 100)
	bt, _ := buf.Pull()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bt.Write([]byte("HelloWorld"))
		bt.Reset()
	}
	b.ReportAllocs()
}

func TestBuffer_Next(t *testing.T) {
	b := NewPermanentBuffer(make([]byte, 6))
	d0 := b.Next(1)
	copy(d0[:], []byte{1})
	d12 := b.Next(2)
	copy(d12[:], []byte{1, 2})
	d35 := b.Next(3)
	copy(d35[:], []byte{3, 4, 5})
	if b.Bytes()[0] != 1 {
		t.Fatal("wrong 0 byte")
	}
	if b.Bytes()[1] != 1 {
		t.Fatal("wrong 1 byte")
	}
	if b.Bytes()[2] != 2 {
		t.Fatal("wrong 2 byte")
	}
	if b.Bytes()[3] != 3 {
		t.Fatal("wrong 3 byte")
	}
	if b.Bytes()[4] != 4 {
		t.Fatal("wrong 3 byte")
	}
	if b.Bytes()[5] != 5 {
		t.Fatal("wrong 5 byte")
	}
	b.Reset()
	b.Write([]byte{1, 2, 3, 4, 5, 6})
	if b.Bytes()[0] != 1 {
		t.Fatal("wrong 0 byte")
	}
	if b.Bytes()[1] != 2 {
		t.Fatal("wrong 1 byte")
	}
	if b.Bytes()[2] != 3 {
		t.Fatal("wrong 2 byte")
	}
	if b.Bytes()[3] != 4 {
		t.Fatal("wrong 3 byte")
	}
	if b.Bytes()[4] != 5 {
		t.Fatal("wrong 3 byte")
	}
	if b.Bytes()[5] != 6 {
		t.Fatal("wrong 5 byte")
	}
}
