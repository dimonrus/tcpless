package tcpless

import (
	"fmt"
	"github.com/dimonrus/gocli"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func getTLSServerConfig() *Config {
	config := getTestConfig()
	config.TLS = TLSConfig{
		Enabled:  true,
		CertPath: "resource/server.crt",
		KeyPath:  "resource/server.pem",
	}
	return config
}

func getTLSClientConfig() *Config {
	config := getTestConfig()
	config.TLS = TLSConfig{
		Enabled:  true,
		CertPath: "resource/GetFreeClient.crt",
		KeyPath:  "resource/GetFreeClient.pem",
		CaPath:   "resource/ca.crt",
	}
	return config
}

func HelloTLS(client IClient) {
	atomic.AddInt32(&rps, 1)
	entity := TestUser{}
	err := client.Parse(&entity)
	if err != nil {
		fmt.Println(err)
	}
	resp := TestOkResponse{Msg: "ok"}
	err = client.Ask("response", resp)
	if err != nil {
		fmt.Println(err)
	}
}

func MyTLSHandler(handler Handler) Handler {
	return handler.Reg("Hello", HelloTLS)
}

func TestServer_TLSStart(t *testing.T) {
	server := NewServer(
		getTLSServerConfig(),
		MyTLSHandler(nil),
		NewGobClient,
		gocli.NewLogger(gocli.LoggerConfig{}),
	)
	err := server.Start()
	if err != nil {
		t.Fatal(err)
	}
	go resetRps()
	time.Sleep(time.Second * 3600)
}

func TestTLSClient(t *testing.T) {
	config := getTLSClientConfig()

	requests := 1_000_000
	parallel := 4

	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			client := NewGobClient(config, nil)
			_, err := client.Dial()
			if err != nil {
				fmt.Println(err)
			}
			resp := TestOkResponse{}
			for j := 0; j < requests; j++ {
				err = client.Ask("Hello", getTestUser())
				if err != nil {
					fmt.Println(err)
				}
				err = client.Parse(&resp)
				if err != nil {
					fmt.Println(err)
				}
				if resp.Msg != "ok" {
					fmt.Println("wrong response")
				}
			}
			_ = client.Stream().Release()
		}()
	}
	wg.Wait()
}

func TestTLSConfig_LoadTLSConfig(t *testing.T) {
	config := getTLSClientConfig()
	for i := 0; i < 200; i++ {
		go func() {
			config.TLS.LoadTLSConfig()
		}()
	}
}
