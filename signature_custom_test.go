package tcpless

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"
)

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

	server, client := getTestConnection()

	gClient := NewGobClient()
	gClient.RegisterType(&UserResp{})
	gClient.SetStream(server)

	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(resp)
	if err != nil {
		t.Fatal(err)
	}
	sig.data = b.Bytes()

	buf, index := testBuffer.Pull()
	defer testBuffer.Release(index)

	go client.Write(sig.Encode(buf))
	
	u := &UserResp{}
	res := Response{
		Data: u,
	}
	err = gClient.Parse(&res)
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
