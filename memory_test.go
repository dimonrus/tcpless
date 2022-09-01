package tcpless

import (
	"github.com/dimonrus/gocli"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func MemoryCheck(client IClient) {
	atomic.AddInt32(&rps, 1)
	client.Read()
}

func MemoryHandler(handler Handler) Handler {
	return handler.
		Reg("memory", MemoryCheck)
}

func TestServerMemory(t *testing.T) {
	config := getTestConfig()
	server := NewServer(config, MemoryHandler(nil), NewGobClient, gocli.NewLogger(gocli.LoggerConfig{}))
	err := server.Start()
	if err != nil {
		t.Fatal(err)
	}
	go resetRps()
	time.Sleep(time.Second * 20)
}

func TestClientMemory(t *testing.T) {
	config := getTestConfig()
	requests := 1000000
	parallel := 3
	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			client := NewGobClient(config, nil)
			_, err := client.Dial()
			if err != nil {
				t.Fatal(err)
			}
			for j := 0; j < requests; j++ {
				//time.Sleep(time.Millisecond * 333)
				err = client.Ask("memory", []byte("how about memory"))
				if err != nil {
					t.Fatal(err)
				}
			}
			_ = client.Close()
		}()
	}
	wg.Wait()
}
