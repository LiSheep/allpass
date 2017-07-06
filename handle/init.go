package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/LiSheep/allpass/model"
	"github.com/LiSheep/allpass/helper"
	"os"
)

type handlerConfig struct {
	Pem struct{
		PubKeyPath string	`json:"pub_key_path"`
		PrvKeyPath string	`json:"prv_key_path"`
	}
}

var config handlerConfig
var cret *helper.Credentials

func Initialize(addr string) {
	err := model.Initialize()
	if err != nil {
		panic(err)
	}
	err = helper.Config(&config)
	if err != nil {
		panic(err)
		return
	}
	cret, err = helper.NewCredentials(config.Pem.PubKeyPath, config.Pem.PrvKeyPath)
	if err != nil {
		panic(err)
	}
	router := gin.New()
	gin.DefaultWriter, _ = os.OpenFile("/root/allpass.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Static("/js", "./static/js")
	router.Static("/css", "./static/css")
	router.StaticFile("/favicon.ico", "./static/favicon.ico")
	router.Static("/plugin", "./static/plugin")

	router.StaticFile("/password", "./views/password/list.html")
	router.StaticFile("/password/list", "./views/password/list.html")

	initUserHandler(router)
	initPasswordHandler(router)
	router.Run(addr)
}
