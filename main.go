package main

import (
	"log"
	"water_proccesing/api"
	conf "water_proccesing/config"
)

func main() {
	config := conf.GetConfig()
	app := &api.App{}
	app.Initialize(config)
	log.Println(" - Server listen on port 3000")
	app.Run(":3000")
}
