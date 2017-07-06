package main

import (
	"os"
	"github.com/LiSheep/allpass/helper"
	"github.com/LiSheep/allpass/handle"
	"github.com/LiSheep/allpass/model"
)

type listenerConfig struct{
	Listen struct {
		Address string `json:"address"`
	}
}

func main() {
	configFile := "./config.json"
	args := os.Args
	if args != nil && len(args) == 2 {
		configFile = args[1]
	}
	fmt.Println("Loading config: ", configFile)
	helper.ConfigInitialize(configFile)
	model.Initialize()
	var config listenerConfig
	err := helper.Config(&config)
	if err != nil {
		return
	}
	handle.Initialize(config.Listen.Address)
}
