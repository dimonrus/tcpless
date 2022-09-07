package tcpless

import (
	"fmt"
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
	config := getTestConfig()
	c := NewGobClient(config, nil)
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

// goos: darwin
// goarch: amd64
// pkg: github.com/dimonrus/tcpless
// cpu: Intel(R) Core(TM) i5-8279U CPU @ 2.40GHz
// BenchmarkGobClient_Parse
// BenchmarkGobClient_Parse-8   	  339273	      3188 ns/op	     120 B/op	       5 allocs/op
func BenchmarkGobClient_Parse(b *testing.B) {
	client, server := getTestClientServer()

	user := getTestUser()
	go func(cl IClient) {
		for i := 0; i < b.N; i++ {
			err := cl.Ask("Hello", user)
			if err != nil {
				fmt.Println(err)
			}
		}
	}(client)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userModel := TestUser{}
		err := server.Parse(&userModel)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
}
