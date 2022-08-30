package tcpless

import (
	"github.com/dimonrus/gocli"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func MemoryCheck(client IClient) {
	atomic.AddInt32(&rps, 1)
	client.Stream().Buffer().Reset()
}

func MemoryHandler(handler Handler) Handler {
	return handler.
		Reg("memory", MemoryCheck)
}

func TestServerMemory(t *testing.T) {
	config := Config{
		Address: &net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 900,
		},
		Limits: ConnectionLimit{
			MaxConnections:   5,
			SharedBufferSize: 1024,
			MaxIdle:          time.Second * 10,
		},
	}
	server := NewServer(config, MemoryHandler(nil), NewGobClient, gocli.NewLogger(gocli.LoggerConfig{}))
	err := server.Start()
	if err != nil {
		t.Fatal(err)
	}
	go resetRps()
	time.Sleep(time.Second * 20)
}

func TestClientMemory(t *testing.T) {
	address := &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 900,
	}
	requests := 1000000
	parallel := 5
	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			client := NewGobClient(nil)
			err := client.Dial(address)
			if err != nil {
				t.Fatal(err)
			}
			//var response *GobSignature
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
