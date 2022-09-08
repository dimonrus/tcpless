package main

import (
	"fmt"
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
	cc, err := tcpless.ConcurrentClient(int(config.Limits.MaxConnections), 0, NewHelloClient, config, app.GetLogger())
	if err != nil {
		app.FailMessage(err.Error())
	}
	cc.RegisterType(&TestUserUserCreate{})

	request := make(chan any, cc.GetConcurrent())
	result := make(chan TestUserUserCreate, 1_000_000)
	wg := sync.WaitGroup{}
	wg.Add(1_000_000)

	cc.Call("api.v1.hello", request, func(client tcpless.IClient) {
		defer wg.Done()
		resp := TestResponse{
			Data: &TestUserUserCreate{},
		}
		err = client.Parse(&resp)
		if err != nil {
			client.Logger().Errorln(err)
		}
		result <- *resp.Data.(*TestUserUserCreate)
	})

	for i := 0; i < 1000000; i++ {
		request <- getTestUser()
	}
	close(request)
	wg.Wait()
	app.SuccessMessage(fmt.Sprintf("received %v results", len(result)))
}
