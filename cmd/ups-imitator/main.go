package main

import (
	"log"

	"github.com/alex11prog/ups-imitator/internal/apiserver"
	"github.com/alex11prog/ups-imitator/internal/app/imitator"
	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/goburrow/modbus"
)

const (
	confPath        = "conf/config.toml"
)

//	@title			UPS-имитатор - OpenAPI спецификация
//	@version		v1.0.0
//	@description	REST API для UPS имитатора

// host		localhost:8080
// @BasePath	/
func main() {
	conf := model.NewConfig(confPath)

	var clients []modbus.Client

	// init tcp modbus clients
	for _, upsAddr := range conf.UpsAddresses {
		handler := modbus.NewTCPClientHandler(upsAddr)
		if err := handler.Connect(); err != nil {
			log.Fatal(err)
		}
		defer handler.Close()
		clients = append(clients, modbus.NewClient(handler))
	}

	imitator := imitator.New(clients, conf)
	imitator.Start()

	if err := apiserver.StartServer(conf.RestApiBindAddr, imitator); err != nil {
		log.Fatal("apiserver.StartServer error! ", err)
	}
}
