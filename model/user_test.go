package model

import (
	"testing"
	"github.com/LiSheep/allpass/helper"
	"sync"
)

var cret *helper.Credentials
func InitTestEnv(t *testing.T) {
	err := helper.ConfigInitialize("../config_dev.json")
	if err != nil {
		t.Error(err.Error())
	}

	cret, err = helper.NewCredentials("../keys/app.rsa.pub", "../keys/app.rsa")
	if err != nil {
		t.Error(err.Error())
	}
	err = Initialize()
	if err != nil {
		t.Error(err.Error())
	}
}
func TestUserAdd_Validate(t *testing.T) {
	err := helper.ConfigInitialize("../config.json")
	if err != nil {
		t.Error(err.Error())
	}
	cret, err = helper.NewCredentials("../keys/app.rsa.pub", "../keys/app.rsa")
	if err != nil {
		t.Error(err.Error())
	}
	err = Initialize()
	if err != nil {
		t.Error(err.Error())
	}
	err = Initialize()
	if err != nil {
		t.Error(err.Error())
	}
	m.db.C("user").RemoveAll(nil)
	s_pass := cret.EncodeHmacSha512([]byte("123"))
	err = UserAdd("litc", s_pass)
	if err != nil {
		t.Error(err.Error())
	}
	result, err := UserValidate("litc", s_pass)
	if err != nil {
		t.Error(err.Error())
	}
	if !result {
		t.Error("validate fail")
	}
	sync.WaitGroup{}
}

func TestUserDelete(t *testing.T) {
	err := helper.ConfigInitialize("../config.json")
	if err != nil {
		t.Error(err.Error())
	}
	cret, err = helper.NewCredentials("../keys/app.rsa.pub", "../keys/app.rsa")
	if err != nil {
		t.Error(err.Error())
	}
	err = Initialize()
	if err != nil {
		t.Error(err.Error())
	}
	m.db.C("user").RemoveAll(nil)
	s_pass := cret.EncodeHmacSha512([]byte("123"))
	UserAdd("litc", s_pass)
	err = UserDelete("litc")
	if err != nil {
		t.Error(err.Error())
	}
	n, err := m.db.C("user").Count()
	if err != nil {
		t.Error(err.Error())
	}
	if n != 0 {
		t.Error("UserDelete fail")
	}
}

func TestUserFind(t *testing.T) {
	err := helper.ConfigInitialize("../config.json")
	if err != nil {
		t.Error(err.Error())
	}
	cret, err = helper.NewCredentials("../keys/app.rsa.pub", "../keys/app.rsa")
	if err != nil {
		t.Error(err.Error())
	}
	err = Initialize()
	if err != nil {
		t.Error(err.Error())
	}
	m.db.C("user").RemoveAll(nil)
	s_pass := cret.EncodeHmacSha512([]byte("123"))
	UserAdd("litc", s_pass)
	u, err := UserFind("litc")
	if err != nil {
		t.Error(err.Error())
	}
	if string(u.Pass) != string(s_pass) {
		t.Error("FindUser error")
	}
}

func TestUserPassUpdate(t *testing.T) {
	err := helper.ConfigInitialize("../config.json")
	if err != nil {
		t.Error(err.Error())
	}
	cret, err = helper.NewCredentials("../keys/app.rsa.pub", "../keys/app.rsa")
	if err != nil {
		t.Error(err.Error())
	}
	err = Initialize()
	if err != nil {
		t.Error(err.Error())
	}
	m.db.C("user").RemoveAll(nil)
	s_pass := cret.EncodeHmacSha512([]byte("123"))
	UserAdd("litc", s_pass)

	s_pass = cret.EncodeHmacSha512([]byte("321"))
	err = UserPassUpdate("litc", s_pass)
	if err != nil {
		t.Error(err.Error())
	}
	u, err := UserFind("litc")
	if err != nil {
		t.Error(err.Error())
	}
	if string(u.Pass) != string(s_pass) {
		t.Error("UserPassUpdate fail")
	}
}