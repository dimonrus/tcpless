package tcpless

import (
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

func TestEncodeDecodeMultipleSend(t *testing.T) {
	client, server := getTestClientServer()

	go func(cl IClient) {
		for i := 0; i < 5000; i++ {
			user := getTestUser()
			err := cl.Ask("Hello", user)
			if err != nil {
				t.Fatal(err)
			}
		}
	}(client)

	for i := 0; i < 3000; i++ {
		user := &TestUser{}
		err := server.Parse(user)
		if err != nil {
			t.Fatal(err)
		}
		if user.Number != 455000 {
			t.Fatal("wrong encode decode number")
		}
		if *user.Name != "ДобрыйДень" {
			t.Fatal("wrong encode decode name")
		}
	}

	for i := 0; i < 2000; i++ {
		user := &TestUser{}
		_, err := server.Read()
		if err != nil {
			t.Fatal(err)
		}
		err = server.Parse(user)
		if err != nil {
			t.Fatal(err)
		}
		if user.Number != 455000 {
			t.Fatal("wrong encode decode number")
		}
		if *user.Name != "ДобрыйДень" {
			t.Fatal("wrong encode decode name")
		}
	}
}
