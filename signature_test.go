package tcpless

import (
	"bytes"
	"net"
	"testing"
)

var testBuffer = CreateBuffer(10, 100)

func getTestConnection() (server Connection, client net.Conn) {
	conn := &connection{done: make(chan struct{})}
	conn.Conn, client = net.Pipe()
	conn.buffer, conn.index = testBuffer.Pull()
	return conn, client
}

func TestGobSignature_Encode(t *testing.T) {
	sig := GobSignature{route: []byte("Hello"), data: []byte("HelloWorld")}
	buf, index := testBuffer.Pull()
	defer testBuffer.Release(index)
	for i := 0; i < 2000; i++ {
		data := sig.Encode(buf)
		reader := bytes.NewBuffer(data)
		h, _ := sig.Decode(reader, buf)
		if h.Len() != 10 || sig.Route() != "Hello" || string(sig.Data()) != "HelloWorld" {
			t.Fatal("wrong encode decode")
		}
	}
}

func TestGobSignature_Decode(t *testing.T) {
	data := []byte{5, 1, 10, 72, 101, 108, 108, 111, 72, 101, 108, 108, 111, 87, 111, 114, 108, 100}
	reader := bytes.NewBuffer(data)
	buf, index := testBuffer.Pull()
	defer testBuffer.Release(index)
	sig := GobSignature{}
	h, err := sig.Decode(reader, buf)
	if err != nil {
		t.Fatal(err)
	}
	if h.Route() != "Hello" {
		t.Fatal("wrong decode route")
	}
	if string(h.Data()) != "HelloWorld" {
		t.Fatal("wrong decode data")
	}
}

func BenchmarkGobSignature_Encode(b *testing.B) {
	sig := GobSignature{route: []byte("Hello"), data: []byte("HelloWorld")}
	buf, index := testBuffer.Pull()
	defer testBuffer.Release(index)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sig.Encode(buf)
	}
	b.ReportAllocs()
}

func BenchmarkGobSignature_Decode(b *testing.B) {
	data := []byte{5, 1, 10, 72, 101, 108, 108, 111, 72, 101, 108, 108, 111, 87, 111, 114, 108, 100}
	reader := bytes.NewBuffer(data)
	sig := GobSignature{}
	buf, _ := testBuffer.Pull()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader.Write(data)
		_, err := sig.Decode(reader, buf)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.ReportAllocs()
}
