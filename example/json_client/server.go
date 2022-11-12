package main

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/tcpless"
	"sync"
)

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
				var err error
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
