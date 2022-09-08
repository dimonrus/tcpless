package main

import (
	"fmt"
	"github.com/dimonrus/tcpless"
	"runtime"
	"sync"
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

	// TestUserUserCreate User create data response
	TestUserUserCreate struct {
		// User Id
		Id *int64
		// Created time
		CreatedAt *time.Time
	}

	// TestResponse Universal response
	TestResponse struct {
		// Message
		Message *string
		// Any data(in case &TestUserUserCreate{})
		Data any
	}
)

var (
	rps          int32
	ticker       = time.NewTicker(time.Millisecond * 1000)
	so           sync.Once
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

func getTestResponse() TestResponse {
	id := int64(1235813)
	now := time.Now()
	c := TestUserUserCreate{
		Id:        &id,
		CreatedAt: &now,
	}
	r := TestResponse{
		Message: new(string),
		Data:    &c,
	}
	*r.Message = "User created successfully."
	return r
}

func HelloHandler(client tcpless.IClient) {
	atomic.AddInt32(&rps, 1)
	entity := TestUser{}
	err := client.Parse(&entity)
	if err != nil {
		panic(err)
	}
	so.Do(func() {
		client.Signature().Encryptor().RegisterType(&TestUserUserCreate{})
	})
	if entity.Number != 455000 {
		panic("json handler does not works")
	}
	err = client.Ask("", getTestResponse())
	if err != nil {
		panic(err)
	}
}

// Handler Custom routes
func Handler(handler tcpless.Handler) tcpless.Handler {
	api := handler.Route("api")
	v1 := api.Sub("v1")
	v1.Handle("hello", HelloHandler)
	return handler
}
