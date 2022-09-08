package main

import (
	"fmt"
	"github.com/dimonrus/tcpless"
	"runtime"
	"sync/atomic"
	"time"
)

type (
	TestUser struct {
		// User name
		Name *string `json:"name"`
		// Some string value
		Some string `json:"some"`
		// Number
		Number int `json:"number"`
	}
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

func getTestUser() TestUser {
	u := TestUser{
		Name:   new(string),
		Some:   "SomeCustomValue",
		Number: 455000,
	}
	*u.Name = "ДобрыйДень"
	return u
}

func HelloHandler(client tcpless.IClient) {
	atomic.AddInt32(&rps, 1)
	entity := TestUser{}
	err := client.Parse(&entity)
	if err != nil {
		panic(err)
	}
}

func V1Hook(client tcpless.IClient) {
	client.Logger().Println("Len of incoming message: ", client.Signature().Len())
}

// Handler Custom routes
func Handler(handler tcpless.Handler) tcpless.Handler {
	// root handler
	api := handler.Route("api")
	// v1 handler and hook
	v1 := api.Sub("v1").Hook(V1Hook)
	// hello handler
	v1.Handle("hello", HelloHandler)
	return handler
}
