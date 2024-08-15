package model

import (
	"testing"
	"time"
)

func TestConfig(t *testing.T) *Config {
	return &Config{
		UpsSyncInterval: time.Second,
	}
}

func TestUpsParams(t *testing.T) UpsParams {
	return UpsParams{
		InputAcVoltage:       220,
		InputAcCurrent:       5,
		BatGroupVoltage:      54,
		BatGroupCurrent:      0,
		LoadCurrent:          20,
		BatCapacity:          50,
		RemainingBatCapacity: 50,
		SOC:                  1,
		Batteries: [4]BatteryParams{
			{
				Voltage: 13.5,
				Temp:    23.2,
				Resist:  5,
			},
			{
				Voltage: 12.5,
				Temp:    24.4,
				Resist:  5.5,
			},
			{
				Voltage: 11.5,
				Temp:    24,
				Resist:  5.2,
			},
			{
				Voltage: 12.5,
				Temp:    23.5,
				Resist:  5.1,
			},
		},
	}
}
