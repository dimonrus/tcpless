package tcpless

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	rps          int32
	ticker       = time.NewTicker(time.Millisecond * 1000)
	m            runtime.MemStats
	memoryReport = map[string]uint64{
		"allocated":          0,
		"total_allocated":    0,
		"system":             0,
		"garbage_collectors": 0}
)

type TestOkResponse struct {
	Msg string
}

func getTestOkResponse() TestOkResponse {
	return TestOkResponse{Msg: "ok"}
}

func printMemStat() {
	runtime.ReadMemStats(&m)
	memoryReport["allocated"] = m.Alloc
	memoryReport["total_allocated"] = m.TotalAlloc
	memoryReport["system"] = m.Sys
	memoryReport["garbage_collectors"] = uint64(m.NumGC)
	fmt.Println(memoryReport)
}

func resetRps() {
	for range ticker.C {
		fmt.Println("rps is: ", atomic.LoadInt32(&rps))
		atomic.StoreInt32(&rps, 0)
		printMemStat()
	}
}

var so = &sync.Once{}

func Hello(client IClient) {
	atomic.AddInt32(&rps, 1)
	entity := TestUser{}
	err := client.Parse(&entity)
	if err != nil {
		fmt.Println(err)
	}
	resp := TestOkResponse{Msg: "ok"}
	err = client.Ask("response", resp)
	if err != nil {
		fmt.Println(err)
	}
}

func MyHandler(handler Handler) Handler {
	return handler.
		Reg("Hello", Hello)
}

func TestServer(t *testing.T) {
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
	server := NewServer(config, MyHandler(nil), NewGobClient, gocli.NewLogger(gocli.LoggerConfig{}))
	err := server.Start()
	if err != nil {
		t.Fatal(err)
	}
	go resetRps()
	time.Sleep(time.Second * 20)
	//c := make(chan os.Signal)
	//<-c
}

func TestClient(t *testing.T) {
	address := &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 900,
	}

	requests := 1000000
	parallel := 2

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
			resp := TestOkResponse{}
			for j := 0; j < requests; j++ {
				err = client.Ask("Hello", getTestUser())
				if err != nil {
					t.Fatal(err)
				}
				err = client.Parse(&resp)
				if err != nil {
					t.Fatal(err)
				}
				if resp.Msg != "ok" {
					t.Fatal("wrong response")
				}
			}
			_ = client.Close()
		}()
	}
	wg.Wait()
}
