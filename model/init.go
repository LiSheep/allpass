package model

import (
	"gopkg.in/mgo.v2"
	"github.com/LiSheep/allpass/helper"
	"fmt"
	"time"
)

type ModelConfig struct {
	DB struct {
		Host string
		Name string
		User string
		Pass string
	}
}

type _model struct {
	session *mgo.Session
	db *mgo.Database
}

var config ModelConfig
var m _model

func Initialize() (err error) {
	err = helper.Config(&config)
	if err != nil {
		panic(err)
		return
	}
	err = initConnection(config.DB.Host, config.DB.Name, config.DB.User, config.DB.Pass)
	if err != nil {
		panic(err)
		return
	}
	return
}

func initConnection(host, dbName, user, pass string) (err error) {
	m.session, err = mgo.Dial("mongodb://"+host+"/"+dbName)
	if err != nil {
		return err
	}
	err = m.session.DB(dbName).Login(user, pass)
	if err != nil {
		return err
	}
	m.db = m.session.DB("")
	go keepAlive(host, dbName, user, pass)
	return nil
}

func keepAlive(host, dbName, user, pass string) {
	for {
		var err error
		if m.session != nil {
			err = m.session.Ping()
		}
		if m.session == nil || err != nil {
			fmt.Println("connect mongo fail")
			m.session, err = mgo.Dial("mongodb://"+host+"/"+dbName)
			if err != nil {
				m.session = nil
				fmt.Println(err)
				continue
			}
			err = m.session.DB(dbName).Login(user, pass)
			if err != nil {
				m.session = nil
				fmt.Println(err)
				continue
			}
			m.db = m.session.DB("")
		}
		time.Sleep(5*time.Second)
	}
}
