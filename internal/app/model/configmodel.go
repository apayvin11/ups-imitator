package model

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	validation "github.com/go-ozzo/ozzo-validation"
)

type Config struct {
	UpsAddr         string        `toml:"ups_addr"`
	RestApiBindAddr string        `toml:"rest_api_bind_addr"`
	UpsSyncInterval time.Duration `toml:"ups_sync_interval"` // sec

	CycleChangeTimeout time.Duration `toml:"cycle_change_timeout"` // charge or discharge (sec)

	DefaultInputAcVoltage float32 `toml:"default_input_ac_voltage"` // V
	MaxBatGroupVoltage    float32 `toml:"max_bat_group_voltage"`    // V
	MinBatGroupVoltage    float32 `toml:"min_bat_group_voltage"`    // V
	LoadPower             float32 `toml:"load_power"`               // W
	DefaultBatCapacity    float32 `toml:"default_bat_capacity"`     // Ah
	ChargeCurrentLimit    float32 `toml:"charge_current_limit"`     // A
	LowSocTriggerAlarm    float32 `toml:"low_soc_trigger_alarm"`    // percent
}

func (conf *Config) validate() error {
	return validation.ValidateStruct(
		conf,
		validation.Field(&conf.UpsAddr, validation.Required),
		validation.Field(&conf.RestApiBindAddr, validation.Required),
		validation.Field(&conf.UpsSyncInterval, validation.Required, validation.Min(time.Second)),
		validation.Field(&conf.CycleChangeTimeout, validation.Required, validation.Min(time.Second)),
		validation.Field(&conf.DefaultInputAcVoltage, validation.Required, validation.Min(float32(150)), validation.Max(float32(300))),
		validation.Field(&conf.MaxBatGroupVoltage, validation.Required, validation.Min(float32(52)), validation.Max(float32(100))),
		validation.Field(&conf.MinBatGroupVoltage, validation.Required, validation.Min(float32(12)), validation.Max(float32(50))),
		validation.Field(&conf.LoadPower, validation.Required, validation.Min(float32(100)), validation.Max(float32(200000))),
		validation.Field(&conf.DefaultBatCapacity, validation.Required, validation.Min(float32(10)), validation.Max(float32(1000))),
		validation.Field(&conf.ChargeCurrentLimit, validation.Required, validation.Min(float32(10)), validation.Max(float32(500))),
		validation.Field(&conf.LowSocTriggerAlarm, validation.Required, validation.Max(float32(0.5))),
	)
}

func NewConfig(configPath string) (*Config, error){
	conf := &Config{}
	_, err := toml.DecodeFile(configPath, conf)
	if err != nil {
		return nil, fmt.Errorf("toml decode file config error: %v", err)
	}
	conf.UpsSyncInterval *= time.Second
	conf.CycleChangeTimeout *= time.Second
	if err := conf.validate(); err != nil {
		return nil, err
	}
	return conf, nil
}
