package main

import (
	"github.com/dimonrus/gocli"
	"github.com/dimonrus/tcpless"
	"sync"
)

// StartServer start server
func StartServer(config *tcpless.Config, application gocli.Application) *tcpless.Server {
	server := tcpless.NewServer(config, Handler(nil), NewHelloClient, App.GetLogger())
	err := server.Start()
	if err != nil {
		application.FatalError(err)
	}
	go resetRps()
	return server
}

// StartClient start client and make call
func StartClient(config *tcpless.Config, app gocli.Application) {
	requests := 5
	parallel := int(config.Limits.MaxConnections)

	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			client := NewHelloClient(config, app.GetLogger()).(*HelloClient)
			_, err := client.Dial()
			if err != nil {
				app.FailMessage(err.Error())
				return
			}
			for j := 0; j < requests; j++ {
				err = client.Hello([]byte("hello my friend. How are you?"))
				if err != nil {
					app.FailMessage(err.Error())
					return
				}
				var sig tcpless.ISignature
				sig, err = client.Read()
				if err != nil {
					app.FailMessage(err.Error())
					return
				}
				client.Logger().Println(string(sig.Data()))
			}
			_ = client.Stream().Release()
		}()
	}
	wg.Wait()
}
