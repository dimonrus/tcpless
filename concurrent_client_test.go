package tcpless

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"sync"
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
		client.Signature().Encryptor().RegisterType(&TestUserUserCreate{})
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
	time.Sleep(time.Second * 30)
	//c := make(chan os.Signal)
	//<-c
}
func TestConcurrentClient_Ask(t *testing.T) {
	cc, err := ConcurrentClient(5, 0, NewGobClient, getTestConfig(), gocli.NewLogger(gocli.LoggerConfig{}))
	if err != nil {
		t.Fatal(err)
	}
	cc.RegisterType(&TestUserUserCreate{})

	request := make(chan any, cc.GetConcurrent())
	result := make(chan TestUserUserCreate, 1_000_000)
	wg := sync.WaitGroup{}
	wg.Add(1_000_000)

	cc.Call("api.hello", request, func(client IClient) {
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

	t.Log(len(result))
	t.Log(<-result)
}
