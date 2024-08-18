package main

import (
	"fmt"
	"log"

	"github.com/alex11prog/ups-imitator/internal/apiserver"
	"github.com/alex11prog/ups-imitator/internal/app/imitator"
	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/goburrow/modbus"
)

const (
	confPath = "conf/config.toml"
)

//	@title		UPS-imitator - OpenAPI specification
//	@version	v1.0.0

// host		localhost:8080
//
//	@BasePath	/
func main() {
	conf, err := model.NewConfig(confPath)
	if err != nil {
		log.Fatal(err)
	}

	handler := modbus.NewTCPClientHandler(conf.UpsAddr)
	if err := handler.Connect(); err != nil {
		log.Fatal(err)
	}
	defer handler.Close()
	client := modbus.NewClient(handler)

	imitator := imitator.New(client, conf)
	imitator.Start()
	
	go func() {
		if err := apiserver.StartServer(conf.RestApiBindAddr, imitator); err != nil {
			log.Fatal("apiserver startup error! ", err)
		}
	}()

	fmt.Println("press enter to quit")
	var s string
	fmt.Scanln(&s)
}
