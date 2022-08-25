package tcpless

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/dimonrus/gocli"
	"net"
	"os"
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

type TestUser struct {
	Id        *int64
	Name      *string
	Some      string
	Number    int
	CreatedAt *time.Time
}

func getTestUser() TestUser {
	u := TestUser{
		Id:        new(int64),
		Name:      new(string),
		Some:      "olololo",
		Number:    455555,
		CreatedAt: new(time.Time),
	}
	*u.Id = 1444
	*u.Name = "Boyarskij"
	*u.CreatedAt = time.Now()
	return u
}

type UserResp struct {
	Id *int64
}

type Response struct {
	Message *string
	Data    any
}

func TestSig(t *testing.T) {
	uu := &UserResp{
		Id: new(int64),
	}
	resp := Response{
		Message: new(string),
		Data:    uu,
	}
	*uu.Id = 100
	*resp.Message = "Some messafe"
	sig := GobSignature{route: "some"}
	sig.RegisterType(&UserResp{})
	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(resp)
	if err != nil {
		t.Fatal(err)
	}
	sig.data = b.Bytes()

	b = bytes.NewBuffer(sig.Encode())
	s := GobSignature{}
	err = s.Decode(b)
	if err != nil {
		t.Fatal(err)
	}

	u := &UserResp{}
	res := Response{
		Data: u,
	}
	err = s.Parse(&res)
	if err != nil {
		t.Fatal(err)
	}
	if *res.Data.(*UserResp).Id != *resp.Data.(*UserResp).Id {
		t.Fatal("wrong encode decode id")
	}
	if *res.Message != *resp.Message {
		t.Fatal("wrong encode decode message")
	}
}

func Hello(ctx context.Context, sig Signature) {
	entity := &TestUser{}
	err := sig.Parse(entity)
	if err != nil {
		fmt.Println(err)
	}
	sig.RegisterType(&UserResp{})

	resp := Response{
		Message: new(string),
		Data: &UserResp{
			Id: new(int64),
		},
	}
	*resp.Message = "howdy"
	*resp.Data.(*UserResp).Id = 100

	_, err = sig.Send(resp)
	if err != nil {
		fmt.Println(err)
	}
}

func TestServer(t *testing.T) {
	config := Config{
		Network: "tcp",
		Address: net.TCPAddr{
			IP:   net.IPv4(0, 0, 0, 0),
			Port: 900,
		},
		Limits: ConnectionLimit{
			MaxConnections: 50,
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

	requests := 1
	parallel := 1

	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			conn, err := net.DialTCP("tcp", nil, address)
			if err != nil {
				t.Fatal(err)
			}
			user := getTestUser()
			sig := GobSignature{route: "Hello", stream: conn}
			us := &UserResp{}
			resp := Response{Data: us}
			sig.RegisterType(us)
			var response *GobSignature
			for j := 0; j < requests; j++ {
				response, err = sig.Send(user)
				if err != nil {
					t.Fatal(err)
				}
				err = response.Parse(&resp)
				if err != nil {
					t.Fatal(err)
				}

			}
			conn.Close()
		}()
	}
	wg.Wait()
}
