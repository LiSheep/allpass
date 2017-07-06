package helper

import (
	"os"
	"encoding/json"
	"io/ioutil"
	"errors"
)

type config struct {
	init bool
	path string
}

var c config

func ConfigInitialize(path string) error {
	c.init = true
	c.path = path
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return nil
}

func Config(val interface {}) error {
	if !c.init {
		return errors.New("not init success")
	}
	data, err := ioutil.ReadFile(c.path)
	err = json.Unmarshal(data, val)
	if err != nil {
		return err
	}

	return nil
}

