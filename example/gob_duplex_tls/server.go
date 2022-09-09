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
	requests := 1_000_000
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
			client.Signature().Encryptor().RegisterType(&TestUserUserCreate{})
			for j := 0; j < requests; j++ {
				err = client.Hello(getTestUser())
				if err != nil {
					app.FailMessage(err.Error())
					return
				}
				resp := TestResponse{
					Data: &TestUserUserCreate{},
				}
				err = client.Parse(&resp)
				if err != nil {
					app.FailMessage(err.Error())
					return
				}
				if *resp.Data.(*TestUserUserCreate).Id != 1235813 {
					app.FailMessage("response id is incorrect")
				}
			}
			_ = client.Stream().Release()
		}()
	}
	wg.Wait()
}
