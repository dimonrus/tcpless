package tcpless

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"
)

// Test User struct
type TestUser struct {
	// User name
	Name *string
	// Some string value
	Some string
	// Number
	Number int
}

// User create data response
type TestUserUserCreate struct {
	// User Id
	Id *int64
	// Created time
	CreatedAt *time.Time
}

type TestResponse struct {
	// Message
	Message *string
	// Any data(in case &TestUserUserCreate{})
	Data any
}

func getTestUser() TestUser {
	u := TestUser{
		Name:   new(string),
		Some:   "SomeCustomValue",
		Number: 455000,
	}
	*u.Name = "ДобрыйДень"
	return u
}

func getTestResponse() TestResponse {
	id := int64(1235813)
	now := time.Now()
	c := TestUserUserCreate{
		Id:        &id,
		CreatedAt: &now,
	}
	r := TestResponse{
		Message: new(string),
		Data:    &c,
	}
	*r.Message = "User created successfully."
	return r
}

func getTestUserGobSignature() Signature {
	sig := &GobSignature{route: []byte("user")}
	buf, index := testBuffer.Pull()
	defer testBuffer.Release(index)
	user := getTestUser()
	b := bytes.NewBuffer(nil)
	err := gob.NewEncoder(b).Encode(user)
	if err != nil {
		panic(err)
	}
	sig.data = b.Bytes()
	sig.Encode(buf)
	return sig
}

func TestEncodeDecode(t *testing.T) {
	server, client := getTestPipe()

	go func(cl Connection) {
		gClient := NewGobClient()
		gClient.SetStream(cl)
		err := gClient.Send("user", getTestUser())
		if err != nil {
			t.Fatal(err)
		}
	}(client)

	sClient := NewGobClient()
	sClient.SetStream(server)
	user := &TestUser{}
	err := sClient.Parse(&user)
	if err != nil {
		t.Fatal(err)
	}

	//err = gClient.Parse(&res)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if *res.Data.(*UserResp).Id != *resp.Data.(*UserResp).Id {
	//	t.Fatal("wrong encode decode id")
	//}
	//if *res.Message != *resp.Message {
	//	t.Fatal("wrong encode decode message")
	//}
}

func BenchmarkClassicGobEncodeDecode(b *testing.B) {
	user := getTestUser()

	bt := make([]byte, 0, 1024)
	buf := bytes.NewBuffer(bt)
	e := gob.NewEncoder(buf)
	d := gob.NewDecoder(buf)

	bt1 := make([]byte, 0, 1024)
	buf1 := bytes.NewBuffer(bt1)
	e1 := gob.NewEncoder(buf1)
	d1 := gob.NewDecoder(buf1)

	sign := TestUser{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := e.Encode(&user)
		if err != nil {
			b.Fatal(err)
		}
		err = d.Decode(&sign)
		if err != nil {
			b.Fatal(err)
		}

		err = e1.Encode(&user)
		if err != nil {
			b.Fatal(err)
		}
		err = d1.Decode(&sign)
		if err != nil {
			b.Fatal(err)
		}
		buf1.Reset()
	}
	b.ReportAllocs()
}
