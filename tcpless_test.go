package tcpless

import (
	"context"
	"fmt"
	"github.com/dimonrus/gocli"
	"net"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func StatusMessage(ctx context.Context, client IClient, sig Signature) {

}

func MyHandler(handler Handler) Handler {
	return handler.
		Reg("Hello", Hello).
		Reg("StatusMessage", StatusMessage)
}

var (
	rps          int32
	ticker       = time.NewTicker(time.Second)
	m            runtime.MemStats
	memoryReport = map[string]uint64{
		"allocated":          0,
		"total_allocated":    0,
		"system":             0,
		"garbage_collectors": 0}
)

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

func Hello(ctx context.Context, client IClient, sig Signature) {
	atomic.AddInt32(&rps, 1)
	entity := TestUser{}
	err := client.Parse(&entity)
	if err != nil {
		fmt.Println(err)
	}
	//sig.RegisterType(&UserResp{})
	//
	//resp := Response{
	//	Message: new(string),
	//	Data: &UserResp{
	//		Id: new(int64),
	//	},
	//}
	//*resp.Message = "howdy"
	//*resp.Data.(*UserResp).Id = 100
	//
	//_, err = sig.Send(resp)
	//if err != nil {
	//	fmt.Println(err)
	//}
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
	//time.Sleep(time.Second * 20)
	c := make(chan os.Signal)
	<-c
}

func TestClient(t *testing.T) {
	address := &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 900,
	}

	requests := 100
	parallel := 2

	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			client := NewGobClient()
			err := client.Dial(address)
			if err != nil {
				t.Fatal(err)
			}
			//us := &UserResp{}
			//resp := Response{Data: us}
			//sig.RegisterType(us)
			//var response *GobSignature
			for j := 0; j < requests; j++ {
				time.Sleep(time.Millisecond * 333)
				err = client.Send("Hello", getTestUser())
				if err != nil {
					t.Fatal(err)
				}

				//err = response.Parse(&resp)
				//if err != nil {
				//	t.Fatal(err)
				//}

			}
			_ = client.Close()
		}()
	}
	wg.Wait()
}
