package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Config_validate(t *testing.T) {
	testCases := []struct {
		name    string
		config  func() *Config
		isValid bool
	}{
		{
			name: "valid",
			config: func() *Config {
				return TestConfig(t)
			},
			isValid: true,
		},
		{
			name: "invalid UpsAddr",
			config: func() *Config {
				conf := TestConfig(t)
				conf.UpsAddr = ""
				return conf
			},
			isValid: false,
		},
		{
			name: "invalid RestApiBindAddr",
			config: func() *Config {
				conf := TestConfig(t)
				conf.RestApiBindAddr = ""
				return conf
			},
			isValid: false,
		},
		{
			name: "invalid UpsSyncInterval",
			config: func() *Config {
				conf := TestConfig(t)
				conf.UpsSyncInterval = time.Millisecond
				return conf
			},
			isValid: false,
		},
		{
			name: "invalid CycleChangeTimeout",
			config: func() *Config {
				conf := TestConfig(t)
				conf.CycleChangeTimeout = time.Millisecond
				return conf
			},
			isValid: false,
		},
		{
			name: "invalid DefaultInputAcVoltage",
			config: func() *Config {
				conf := TestConfig(t)
				conf.DefaultInputAcVoltage = 140
				return conf
			},
			isValid: false,
		},
		{
			name: "invalid MaxBatGroupVoltage",
			config: func() *Config {
				conf := TestConfig(t)
				conf.MaxBatGroupVoltage = 45
				return conf
			},
			isValid: false,
		},
		{
			name: "invalid MinBatGroupVoltage",
			config: func() *Config {
				conf := TestConfig(t)
				conf.MinBatGroupVoltage = 0
				return conf
			},
			isValid: false,
		},
		{
			name: "invalid LoadPower",
			config: func() *Config {
				conf := TestConfig(t)
				conf.LoadPower = 0
				return conf
			},
			isValid: false,
		},
		{
			name: "invalid DefaultBatCapacity",
			config: func() *Config {
				conf := TestConfig(t)
				conf.DefaultBatCapacity = 0
				return conf
			},
			isValid: false,
		},
		{
			name: "invalid ChargeCurrentLimit",
			config: func() *Config {
				conf := TestConfig(t)
				conf.ChargeCurrentLimit = 0
				return conf
			},
			isValid: false,
		},
		{
			name: "invalid LowSocTriggerAlarm",
			config: func() *Config {
				conf := TestConfig(t)
				conf.LowSocTriggerAlarm = 0
				return conf
			},
			isValid: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.config().validate())
			} else {
				assert.Error(t, tc.config().validate())
			}
		})
	}
}
