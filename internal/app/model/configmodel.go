package model

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
	validation "github.com/go-ozzo/ozzo-validation"
)

type Config struct {
	UpsAddresses    []string      `toml:"ups_addresses"`
	RestApiBindAddr string        `toml:"rest_api_bind_addr"`
	UpsSyncInterval  time.Duration `toml:"ups_sync_interval"` // sec
}

func (conf *Config) validate() error {
	return validation.ValidateStruct(
		conf,
		validation.Field(&conf.UpsAddresses, validation.Required),
		validation.Field(&conf.RestApiBindAddr, validation.Required),
		validation.Field(&conf.UpsSyncInterval, validation.Required, validation.Min(time.Second)),
	)
}

func NewConfig(configPath string) *Config {
	conf := &Config{}
	_, err := toml.DecodeFile(configPath, conf)
	if err != nil {
		log.Fatal("toml decode file config error! ", err)
	}
	conf.UpsSyncInterval *= time.Second
	if err := conf.validate(); err != nil {
		log.Fatal(err)
	}
	return conf
}
