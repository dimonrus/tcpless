package tcpless

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func MemoryCheck(client IClient) {
	atomic.AddInt32(&rps, 1)
	atomic.AddInt32(&counter, 1)
	if string(client.Signature().Data()) != "how about memory" {
		panic("wrong")
	}
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
	time.Sleep(time.Second * 200)
}

func TestClientMemory(t *testing.T) {
	config := getTestConfig()
	requests := 20000000
	parallel := 8
	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			client := NewGobClient(config, gocli.NewLogger(gocli.LoggerConfig{}))
			_, err := client.Dial()
			if err != nil {
				fmt.Println(err)
			}
			for j := 0; j < requests; j++ {
				err = client.AskBytes("memory", []byte("how about memory"))
				if err != nil {
					fmt.Println(err)
					return
				}
				time.Sleep(time.Microsecond)
			}
			_ = client.Stream().Release()
		}()
	}
	wg.Wait()
}
