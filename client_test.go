package tcpless

import (
	"testing"
)

func TestGobClient_Read(t *testing.T) {
	server, client := getTestPipe()
	c := NewGobClient()
	c.SetStream(server)
	go func() {
		for i := 0; i < 5; i++ {
			_, err := client.Connection().Write(HelloHelloWorldSignature)
			if err != nil {
				t.Log(err)
			}
		}
	}()
	for i := 0; i < 5; i++ {
		sig, err := c.Read()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(sig)
	}
}

func TestGobClient_Parse(t *testing.T) {
	getTestUser()
	server, client := getTestPipe()
	c := NewGobClient()
	c.SetStream(server)
	go func() {
		for i := 0; i < 5; i++ {
			client.Connection().Write(HelloHelloWorldSignature)
		}
	}()
	for i := 0; i < 5; i++ {
		sig, err := c.Read()
		if err != nil {
			t.Fatal(err)
		}
		//c.Parse()
		t.Log(sig)
	}
}

// goos: darwin
// goarch: amd64
// pkg: github.com/dimonrus/tcpless
// cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
// BenchmarkGobClient_Read
// BenchmarkGobClient_Read-8   	 1000000	      1052 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGobClient_Read(b *testing.B) {
	server, client := getTestPipe()
	c := NewGobClient()
	c.SetStream(server)
	go func() {
		for i := 0; i < b.N; i++ {
			client.Connection().Write(HelloHelloWorldSignature)
		}
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := c.Read()
		if err != nil {
			b.Fatal(err)
		}
	}
	b.ReportAllocs()
}
