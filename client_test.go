package tcpless

import (
	"testing"
)

func TestGobClient_Read(t *testing.T) {
	server, client := getTestConnection()
	c := NewGobClient()
	c.SetStream(server)
	go func() {
		for i := 0; i < 5; i++ {
			client.Write(HelloHelloWorldSignature)
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

// BenchmarkGobClient_Read-8   	  397808	      2730 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGobClient_Read(b *testing.B) {
	server, client := getTestConnection()
	c := NewGobClient()
	c.SetStream(server)
	go func() {
		for i := 0; i < b.N; i++ {
			client.Write(HelloHelloWorldSignature)
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
