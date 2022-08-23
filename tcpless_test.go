package tcpless

import (
	"context"
	"fmt"
	"github.com/dimonrus/gocli"
	"net"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func StatusMessage(ctx context.Context, sig Signature) {

}

func MyHandler(handler Handler) Handler {
	return handler.
		Reg("Hello", Hello).
		Reg("StatusMessage", StatusMessage)
}

var rps int

var ticker = time.NewTicker(time.Second)

//func init() {
//	for range ticker.C {
//		fmt.Println(rps)
//		rps = 0
//	}
//}

func Hello(ctx context.Context, sig Signature) {
	data := make([]byte, sig.Len())
	_, err := sig.Stream().Read(data)
	if err != nil {
		fmt.Println(err)
	}
	//rps++
	//fmt.Println(n)
	fmt.Println(string(data))
	//sig.Stream().Write([]byte("done"))
}

func TestServer(t *testing.T) {
	config := Config{
		Network: "tcp",
		Address: net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 900,
		},
		Limits: ConnectionLimit{
			MaxConnections: 5,
			MaxIdle:        time.Second * 10,
		},
	}
	server := NewServer(config, MyHandler(nil), gocli.NewLogger(gocli.LoggerConfig{}))
	err := server.Start()
	if err != nil {
		t.Fatal(err)
	}
	c := make(chan os.Signal)
	<-c
}

func TestClient(t *testing.T) {
	address := &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 900,
	}

	requests := 1_0
	parallel := 6

	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			conn, err := net.DialTCP("tcp", nil, address)
			if err != nil {
				t.Fatal(err)
			}
			for j := 0; j < requests; j++ {
				message := CreateMessage("Hello", []byte("HelloWorld"+strconv.FormatInt(int64(j), 10)))
				_, err = conn.Write(message)
				if err != nil {
					t.Fatal(err)
				}
			}
			conn.Close()
		}()
	}
	wg.Wait()
}
