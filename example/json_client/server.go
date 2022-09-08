package main

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/tcpless"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type (
	TestOkResponse struct {
		Msg string `json:"msg"`
	}

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

func HelloHandler(client tcpless.IClient) {
	atomic.AddInt32(&rps, 1)
	entity := TestUser{}
	err := client.Parse(&entity)
	if err != nil {
		fmt.Println(err)
	}
	if entity.Number != 455000 {
		panic("json handler does not works")
	}
	resp := TestOkResponse{Msg: "ok"}
	err = client.Ask("", resp)
	if err != nil {
		fmt.Println(err)
	}
}

// JSONHandler Custom routes
func JSONHandler(handler tcpless.Handler) tcpless.Handler {
	api := handler.Route("api")
	api.Handle("hello", HelloHandler)
	return handler
}

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

// StartServer start server
func StartServer(config *tcpless.Config, application gocli.Application) *tcpless.Server {
	server := tcpless.NewServer(config, JSONHandler(nil), NewJSONClient, App.GetLogger())
	err := server.Start()
	if err != nil {
		application.FatalError(err)
	}
	go resetRps()
	return server
}

// StartClient start client and make call
func StartClient(config *tcpless.Config, app gocli.Application) {
	requests := 1_000_000
	parallel := 5

	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			client := NewJSONClient(config, app.GetLogger()).(*JsonClient)
			_, err := client.Dial()
			if err != nil {
				app.FailMessage(err.Error())
				return
			}
			resp := TestOkResponse{}
			for j := 0; j < requests; j++ {
				resp, err = client.Hello(getTestUser())
				if err != nil {
					app.FailMessage(err.Error())
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
