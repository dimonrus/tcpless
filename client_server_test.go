package tcpless

import (
	"fmt"
	"github.com/dimonrus/gocli"
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
	err = client.Ask("", resp)
	if err != nil {
		fmt.Println(err)
	}
}

func MyHandler(handler Handler) Handler {
	api := handler.Route("api")
	api.Hook(func(client IClient) {
		var i int
		_ = i
	})
	api.Handle("Hello", Hello)
	return handler
}

func TestServer(t *testing.T) {
	config := getTestConfig()
	server := NewServer(config, MyHandler(nil), NewGobClient, gocli.NewLogger(gocli.LoggerConfig{}))
	err := server.Start()
	if err != nil {
		t.Fatal(err)
	}
	go resetRps()
	time.Sleep(time.Second * 2)
	server.Stop()
	time.Sleep(time.Second * 2)
	err = server.Start()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 200)
	//c := make(chan os.Signal)
	//<-c
}

func TestClient(t *testing.T) {
	config := getTestConfig()

	requests := 1000000
	parallel := 4

	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			client := NewGobClient(config, nil)
			_, err := client.Dial()
			if err != nil {
				fmt.Println(err)
			}
			resp := TestOkResponse{}
			for j := 0; j < requests; j++ {
				err = client.Ask("api.Hello", getTestUser())
				if err != nil {
					fmt.Println(err)
					return
				}
				err = client.Parse(&resp)
				if err != nil {
					fmt.Println(err)
					return
				}
				if resp.Msg != "ok" {
					fmt.Println("wrong response")
					return
				}
			}
			_ = client.Stream().Release()
		}()
	}
	wg.Wait()
}
