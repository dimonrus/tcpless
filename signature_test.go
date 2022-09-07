package tcpless

import (
	"bytes"
	"net"
	"testing"
)

var testBuffer = CreateBuffer(10, 1024)

var HelloHelloWorldSignature = []byte{5, 0, 0, 10, 72, 101, 108, 108, 111, 72, 101, 108, 108, 111, 87, 111, 114, 108, 100}

func getTestPipe() (server Streamer, client Streamer) {
	srv := &connection{done: make(chan struct{})}
	cl := &connection{done: make(chan struct{})}
	srv.Conn, cl.Conn = net.Pipe()
	srv.buffer, srv.index = testBuffer.Pull()
	cl.buffer, cl.index = testBuffer.Pull()
	return srv, cl
}

func getTestClientServer() (client IClient, server IClient) {
	srv, cl := getTestPipe()
	config := getTestConfig()

	client = NewGobClient(config, nil)
	client.SetStream(cl)

	server = NewGobClient(config, nil)
	server.SetStream(srv)

	return
}

func TestGobSignature_Encode(t *testing.T) {
	sig := Signature{route: []byte("Hello"), data: []byte("HelloWorld")}
	res := Signature{}
	buf, index := testBuffer.Pull()
	defer testBuffer.Release(index)
	reader := bytes.NewBuffer(nil)
	for i := 0; i < 2000000; i++ {
		reader.Write(sig.Encode(buf))
		err := res.Decode(reader, buf)
		if err != nil {
			t.Fatal(err)
		}
		if res.Len() != 10 || sig.Route() != "Hello" || string(sig.Data()) != "HelloWorld" {
			t.Fatal("wrong encode decode")
		}
		reader.Reset()
	}
}

func TestGobSignature_Decode(t *testing.T) {
	reader := bytes.NewBuffer(nil)
	buf, index := testBuffer.Pull()
	defer testBuffer.Release(index)
	sig := Signature{}
	for i := 0; i < 1000; i++ {
		reader.Write(HelloHelloWorldSignature)
		err := sig.Decode(reader, buf)
		if err != nil {
			t.Fatal(err)
		}
		if sig.Route() != "Hello" {
			t.Fatal("wrong decode route")
		}
		if string(sig.Data()) != "HelloWorld" {
			t.Fatal("wrong decode data")
		}
	}
}

// goos: darwin
// goarch: amd64
// pkg: github.com/dimonrus/tcpless
// cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
// BenchmarkGobSignature_Encode
// BenchmarkGobSignature_Encode-8   	65650748	        17.49 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGobSignature_Encode(b *testing.B) {
	sig := Signature{route: []byte("Hello"), data: []byte("HelloWorld")}
	buf, _ := testBuffer.Pull()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sig.Encode(buf)
	}
	b.ReportAllocs()
}

// goos: darwin
// goarch: amd64
// pkg: github.com/dimonrus/tcpless
// cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
// BenchmarkGobSignature_Decode
// BenchmarkGobSignature_Decode-8   	53649286	        21.76 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGobSignature_Decode(b *testing.B) {
	reader := bytes.NewBuffer(nil)
	sig := &Signature{}
	buf, _ := testBuffer.Pull()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader.Write(HelloHelloWorldSignature)
		err := sig.Decode(reader, buf)
		buf.Reset()
		if err != nil {
			b.Fatal(err)
		}
		if string(sig.Data()) != "HelloWorld" {
			b.Fatal("wrong decode")
		}
	}
	b.ReportAllocs()
}
