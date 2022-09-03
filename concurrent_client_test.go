package tcpless

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"sync/atomic"
	"testing"
	"time"
)

func ConcurrentHello(client IClient) {
	atomic.AddInt32(&rps, 1)
	entity := TestUser{}
	err := client.Parse(&entity)
	if err != nil {
		fmt.Println(err)
	}
	so.Do(func() {
		client.RegisterType(&TestUserUserCreate{})
	})
	resp := getTestResponse()
	err = client.Ask("response", resp)
	if err != nil {
		fmt.Println(err)
	}
}

func MyConcurrentHandler(handler Handler) Handler {
	api := handler.Route("api")
	api.Handle("hello", ConcurrentHello)
	return handler
}

func TestConcurrentServer(t *testing.T) {
	config := getTestConfig()
	server := NewServer(config, MyConcurrentHandler(nil), NewGobClient, gocli.NewLogger(gocli.LoggerConfig{}))
	err := server.Start()
	if err != nil {
		t.Fatal(err)
	}
	go resetRps()
	time.Sleep(time.Second * 20)
	//c := make(chan os.Signal)
	//<-c
}

func TestConcurrentClient_Ask(t *testing.T) {
	cc, err := ConcurrentClient(5, 0, NewGobClient, getTestConfig(), gocli.NewLogger(gocli.LoggerConfig{}))
	if err != nil {
		t.Fatal(err)
	}
	user := getTestUser()
	userCreate := TestUserUserCreate{}
	resp := TestResponse{
		Data: &userCreate,
	}
	cc.RegisterType(&userCreate)
	err = cc.Ask("api.hello", user, &resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*resp.Data.(*TestUserUserCreate).Id)
}
