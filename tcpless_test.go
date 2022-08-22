package tcpless

import (
	"context"
	"fmt"
	"github.com/dimonrus/gocli"
	"net"
	"os"
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
	//fmt.Println(string(data))
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

	requests := 1_000_000
	wait := make(chan struct{}, 3)

	go func() {
		conn1, err := net.DialTCP("tcp", nil, address)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < requests; i++ {
			message := CreateMessage("Hello", []byte("HelloWorld"))
			_, err = conn1.Write(message)
			if err != nil {
				t.Fatal(err)
			}
		}
		wait <- struct{}{}
	}()

	go func() {
		conn2, err := net.DialTCP("tcp", nil, address)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < requests; i++ {
			message := CreateMessage("Hello", []byte("HelloWorld"))
			_, err = conn2.Write(message)
			if err != nil {
				t.Fatal(err)
			}
		}
		wait <- struct{}{}
	}()

	go func() {
		conn3, err := net.DialTCP("tcp", nil, address)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < requests; i++ {
			message := CreateMessage("Hello", []byte("HelloWorld"))
			_, err = conn3.Write(message)
			if err != nil {
				t.Fatal(err)
			}
		}
		wait <- struct{}{}
	}()

	go func() {
		conn4, err := net.DialTCP("tcp", nil, address)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < requests; i++ {
			message := CreateMessage("StatusMessage", []byte("StatusMessageExists"))
			_, err = conn4.Write(message)
			if err != nil {
				t.Fatal(err)
			}
		}
		wait <- struct{}{}
	}()

	<-wait
	<-wait
	<-wait
	<-wait
	//resp := make([]byte, 512)
	//n, err := conn.Read(resp)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(string(resp[:n]))
}
