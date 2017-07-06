package model

import (
	"testing"
)

func TestUser_AddPassword(t *testing.T) {
	InitTestEnv(t)
	m.db.C("user").RemoveAll(nil)
	s_pass := cret.EncodeHmacSha512([]byte("123"))
	err := UserAdd("litc", s_pass)
	if err != nil {
		t.Error(err.Error())
	}
	user, err := UserFind("litc")
	if err != nil {
		t.Error(err.Error())
	}
	user.AddPassword("hello", "litc", "123456")
	user.AddPassword("hello2", "litc", "1234567")
	user.RemovePassword("hello")
	user.RemovePassword("hello2")
}

func TestUser_FindPassword(t *testing.T) {
	InitTestEnv(t)
	m.db.C("user").RemoveAll(nil)
	s_pass := cret.EncodeHmacSha512([]byte("123"))
	err := UserAdd("litc", s_pass)
	if err != nil {
		t.Error(err.Error())
	}
	user, err := UserFind("litc")
	if err != nil {
		t.Error(err.Error())
	}
	user.AddPassword("hello", "litc", "123456")
	// need to find again!
	user, err = UserFind("litc")
	if err != nil {
		t.Error(err.Error())
	}
	p := user.FindPassword("hello")
	if p == nil {
		t.Error("findpassword error")
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	InitTestEnv(t)
	m.db.C("user").RemoveAll(nil)
	s_pass := cret.EncodeHmacSha512([]byte("123"))
	err := UserAdd("litc", s_pass)
	if err != nil {
		t.Error(err.Error())
	}
	user, err := UserFind("litc")
	if err != nil {
		t.Error(err.Error())
	}

	user.AddPassword("hello", "litc", "123456")
	user.AddPassword("hahaha", "litc", "123456")

	user, err = UserFind("litc")
	if err != nil {
		t.Error(err.Error())
	}

	updateId := user.FindPassword("hahaha").Id
	err = user.UpdatePassword(updateId, "world", "ltc", "66666")
	if err != nil {
		t.Error(err.Error())
	}

	user, err = UserFind("litc")
	if err != nil {
		t.Error(err.Error())
	}
	if user.FindPassword("hahaha") != nil || user.FindPassword("world") == nil {
		t.Error("update password fail")
	}
}

func TestUser_UpdateAllPasswords(t *testing.T) {
	InitTestEnv(t)
	m.db.C("user").RemoveAll(nil)
	s_pass := cret.EncodeHmacSha512([]byte("123"))
	err := UserAdd("litc", s_pass)
	if err != nil {
		t.Error(err.Error())
	}
	user, err := UserFind("litc")
	if err != nil {
		t.Error(err.Error())
	}

	user.AddPassword("hello", "litc", "123456")
	user.AddPassword("hahaha", "litc", "123456")
	// fresh
	user, err = UserFind("litc")
	if err != nil {
		t.Error(err.Error())
	}
	user.Passwords[0].Secret = "654321"
	user.Passwords[1].Secret = "654321"
	err = user.UpdateAllPasswords(user.Passwords)
	if err != nil {
		t.Error(err.Error())
	}
	// fresh
	user, err = UserFind("litc")
	if err != nil {
		t.Error(err.Error())
	}
	if user.Passwords[0].Secret != "654321" || user.Passwords[1].Secret != "654321" {
		t.Error("update fail")
	}
}

