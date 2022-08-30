package tcpless

import (
	"testing"
)

// goos: darwin
// goarch: amd64
// pkg: github.com/dimonrus/tcpless
// cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
// BenchmarkGobClient_Read
// BenchmarkGobClient_Read-8   	  602586	      1888 ns/op	       0 B/op	       0 allocs/op
func BenchmarkGobClient_Read(b *testing.B) {
	server, client := getTestPipe()
	c := NewGobClient(nil)
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

func BenchmarkGobClient_Parse(b *testing.B) {
	b.Log("hello")
	client, server := getTestClientServer()

	user := getTestUser()
	go func(cl IClient) {
		for i := 0; i < b.N; i++ {
			err := cl.Ask("Hello", user)
			if err != nil {
				b.Fatal(err)
			}
		}
	}(client)

	userModel := &TestUser{}
	for i := 0; i < b.N; i++ {
		err := server.Parse(userModel)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}
