package main

import (
	"fmt"
	"github.com/dimonrus/tcpless"
	"runtime"
	"sync/atomic"
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

func resetRps() {
	for range ticker.C {
		fmt.Println("rps is: ", atomic.LoadInt32(&rps))
		atomic.StoreInt32(&rps, 0)
		printMemStat()
	}
}

func printMemStat() {
	runtime.ReadMemStats(&m)
	memoryReport["allocated"] = m.Alloc
	memoryReport["total_allocated"] = m.TotalAlloc
	memoryReport["system"] = m.Sys
	memoryReport["garbage_collectors"] = uint64(m.NumGC)
	fmt.Println(memoryReport)
}

func HelloHandler(client tcpless.IClient) {
	atomic.AddInt32(&rps, 1)
	client.Logger().Infoln("Message from client: ", string(client.Signature().Data()))
	err := client.AskBytes("", []byte("thank you. i'm fine"))
	if err != nil {
		panic(err)
	}
}

// Handler Custom routes
func Handler(handler tcpless.Handler) tcpless.Handler {
	api := handler.Route("api")
	v1 := api.Sub("v1")
	v1.Handle("exchange", HelloHandler)
	return handler
}
