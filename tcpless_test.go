package tcpless

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/dimonrus/gocli"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func StatusMessage(ctx context.Context, client IClient, sig Signature) {

}

func MyHandler(handler Handler) Handler {
	return handler.
		Reg("Hello", Hello).
		Reg("StatusMessage", StatusMessage)
}

var rps int32

var ticker = time.NewTicker(time.Second)

var idleRps bool

func resetRps() {
	for range ticker.C {
		fmt.Println(atomic.LoadInt32(&rps))
		atomic.StoreInt32(&rps, 0)

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		report := make(map[string]string)
		report["allocated"] = fmt.Sprintf("%v KB", m.Alloc/1024)
		report["total_allocated"] = fmt.Sprintf("%v KB", m.TotalAlloc/1024)
		report["system"] = fmt.Sprintf("%v KB", m.Sys/1024)
		report["garbage_collectors"] = fmt.Sprintf("%v", m.NumGC)
		fmt.Println(report)

	}
}

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
	sig := GobSignature{route: []byte("some")}

	client := GobClient{}
	client.RegisterType(&UserResp{})

	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(resp)
	if err != nil {
		t.Fatal(err)
	}
	sig.data = b.Bytes()

	buf, index := testBuffer.Pull()
	defer testBuffer.Release(index)

	reader := bytes.NewBuffer(sig.Encode(buf))

	s, err := GobSignature{}.Decode(reader, buf)
	if err != nil {
		t.Fatal(err)
	}

	u := &UserResp{}
	res := Response{
		Data: u,
	}
	err = client.Parse(s, &res)
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

func Hello(ctx context.Context, client IClient, sig Signature) {
	if !idleRps {
		go resetRps()
	}
	atomic.AddInt32(&rps, 1)
	entity := &TestUser{}
	err := client.Parse(sig, entity)
	if err != nil {
		fmt.Println(err)
	}
	//sig.RegisterType(&UserResp{})
	//
	//resp := Response{
	//	Message: new(string),
	//	Data: &UserResp{
	//		Id: new(int64),
	//	},
	//}
	//*resp.Message = "howdy"
	//*resp.Data.(*UserResp).Id = 100
	//
	//_, err = sig.Send(resp)
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func TestServer(t *testing.T) {
	config := Config{
		Address: &net.TCPAddr{
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
	//c := make(chan os.Signal)
	//<-c
	time.Sleep(time.Second * 10)
}

func TestClient(t *testing.T) {
	address := &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 900,
	}

	requests := 1_000_000
	parallel := 2

	wg := sync.WaitGroup{}
	wg.Add(parallel)
	for i := 0; i < parallel; i++ {
		go func() {
			defer wg.Done()
			client := GobClient{}
			err := client.Dial(address)
			if err != nil {
				t.Fatal(err)
			}
			//us := &UserResp{}
			//resp := Response{Data: us}
			//sig.RegisterType(us)
			//var response *GobSignature
			for j := 0; j < requests; j++ {
				err = client.Send("Hello", getTestUser())
				if err != nil {
					t.Fatal(err)
				}
				//err = response.Parse(&resp)
				//if err != nil {
				//	t.Fatal(err)
				//}

			}
			_ = client.Close()
		}()
	}
	wg.Wait()
}
