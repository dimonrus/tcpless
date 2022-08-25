package tcpless

import (
	"bytes"
	"testing"
)

func TestGobSignature_Encode(t *testing.T) {
	sig := GobSignature{route: "Hello", data: []byte("HelloWorld")}
	data := sig.Encode()
	b := bytes.NewBuffer(data)
	_ = sig.Decode(b)
	if sig.Len() != 10 {
		t.Fatal("wrong encode decode")
	}
}

func TestGobSignature_Decode(t *testing.T) {
	data := []byte{5, 1, 10, 72, 101, 108, 108, 111, 72, 101, 108, 108, 111, 87, 111, 114, 108, 100}
	buf := bytes.NewBuffer(data)
	sig := GobSignature{}
	err := sig.Decode(buf)
	if err != nil {
		t.Fatal(err)
	}
	if sig.route != "Hello" {
		t.Fatal("wrong decode route")
	}
	if string(sig.data) != "HelloWorld" {
		t.Fatal("wrong decode data")
	}
}

func BenchmarkGobSignature_Encode(b *testing.B) {
	sig := GobSignature{route: "Hello", data: []byte("HelloWorld")}
	for i := 0; i < b.N; i++ {
		_ = sig.Encode()
	}
	b.ReportAllocs()
}

func BenchmarkGobSignature_Decode(b *testing.B) {
	data := []byte{5, 1, 10, 72, 101, 108, 108, 111, 72, 101, 108, 108, 111, 87, 111, 114, 108, 100}
	buf := bytes.NewBuffer(data)
	sig := GobSignature{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sig.Decode(buf)
	}
	b.ReportAllocs()
}
