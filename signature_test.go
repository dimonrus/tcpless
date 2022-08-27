package tcpless

import (
	"bytes"
	"encoding/gob"
	"net"
	"testing"
	"time"
)

var testBuffer = CreateBuffer(10, 100)

func getTestConnection() (server Connection, client net.Conn) {
	conn := &connection{done: make(chan struct{})}
	conn.Conn, client = net.Pipe()
	conn.buffer, conn.index = testBuffer.Pull()
	return conn, client
}

func TestGobSignature_Encode(t *testing.T) {
	sig := GobSignature{route: []byte("Hello"), data: []byte("HelloWorld")}
	res := GobSignature{}
	buf, index := testBuffer.Pull()
	defer testBuffer.Release(index)
	for i := 0; i < 2000000; i++ {
		data := sig.Encode(buf)
		reader := bytes.NewBuffer(data)
		err := res.Decode(reader, buf)
		if err != nil {
			t.Fatal(err)
		}
		if res.Len() != 10 || sig.Route() != "Hello" || string(sig.Data()) != "HelloWorld" {
			t.Fatal("wrong encode decode")
		}
	}
}

func TestGobSignature_Decode(t *testing.T) {
	data := []byte{5, 1, 10, 72, 101, 108, 108, 111, 72, 101, 108, 108, 111, 87, 111, 114, 108, 100}
	reader := bytes.NewBuffer(nil)
	buf, index := testBuffer.Pull()
	defer testBuffer.Release(index)
	sig := GobSignature{}
	for i := 0; i < 1000; i++ {
		reader.Write(data)
		err := sig.Decode(reader, buf)
		if err != nil {
			t.Fatal(err)
		}
		if sig.Route() != "Hello" {
			t.Fatal("wrong decode route")
		}
		if string(sig.Data()) != "HelloWorld" {
			t.Fatal("wrong decode data")
		}
	}
}

func BenchmarkGobSignature_Encode(b *testing.B) {
	sig := GobSignature{route: []byte("Hello"), data: []byte("HelloWorld")}
	buf, _ := testBuffer.Pull()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sig.Encode(buf)
	}
	b.ReportAllocs()
}

func BenchmarkGobSignature_Decode(b *testing.B) {
	data := []byte{5, 1, 10, 72, 101, 108, 108, 111, 72, 101, 108, 108, 111, 87, 111, 114, 108, 100}
	reader := bytes.NewBuffer(data)
	sig := GobSignature{}
	buf, _ := testBuffer.Pull()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader.Write(data)
		err := sig.Decode(reader, buf)
		if err != nil {
			b.Fatal(err)
		}
		if string(sig.Data()) != "HelloWorld" {
			b.Fatal("wrong decode")
		}
	}
	b.ReportAllocs()
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
	s := &GobSignature{}
	err = s.Decode(reader, buf)
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
