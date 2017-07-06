package helper

import (
	"testing"
	"fmt"
)

func TestEncodeAndDecode(t *testing.T) {
	var d []byte
	d = []byte("hello world")
	c, err := NewCredentials("../keys/app.rsa.pub", "../keys/app.rsa")
	if err != nil {
		t.Errorf("err", err)
		return
	}
	secData, err := c.Encode(d)
	if err != nil {
		t.Errorf("err", err)
		return
	}
	result, err := c.Decode(secData)
	if err != nil {
		t.Errorf("err", err)
		return
	}
	if (string(d) != string(result)) {
		t.Errorf("result error: expect %s get %s", string(d), string(result))
	}
}

func TestSHA512(t *testing.T) {
	c, err := NewCredentials("../keys/app.rsa.pub", "../keys/app.rsa")
	if err != nil {
		t.Errorf("err", err)
		return
	}
	d := []byte("hello world")
	result := c.EncodeSHA512(d)
	fmt.Println(string(result))
}

func TestSHA256(t *testing.T) {
	c, err := NewCredentials("../keys/app.rsa.pub", "../keys/app.rsa")
	if err != nil {
		t.Errorf("err", err)
		return
	}
	d := []byte("hello world")
	result := c.EncodeSHA256(d)
	fmt.Println(string(result), len(result))
}

func TestHmacSha512(t *testing.T) {
	c, err := NewCredentials("../keys/app.rsa.pub", "../keys/app.rsa")
	if err != nil {
		t.Errorf("err", err)
		return
	}
	d := []byte("hello world")
	result := c.EncodeHmacSha512(d)
	fmt.Println(string(result), len(result))
}

func TestHmacSha256(t *testing.T) {
	c, err := NewCredentials("../keys/app.rsa.pub", "../keys/app.rsa")
	if err != nil {
		t.Errorf("err", err)
		return
	}
	d := []byte("hello world")
	result := c.EncodeHmacSha256(d)
	fmt.Println(string(result), len(result))
}
