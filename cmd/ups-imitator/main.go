package main

import (
	"log"

	"github.com/alex11prog/ups-imitator/internal/app/imitator"
	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/goburrow/modbus"
)

const confPath = "conf/config.toml"

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
}
