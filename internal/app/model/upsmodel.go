package model

type BatteryParams struct {
	Voltage float32 `json:"voltage" example:"12"`
	Temp    float32 `json:"temp" example:"24"`
	Resist  float32 `json:"resist" example:"5"`
}

type Alarms struct {
	UpcInBatteryMode bool `json:"upc_in_battery_mode" example:"false"`
	LowBattery       bool `json:"low_battery" example:"false"`
	Overload         bool `json:"overload" example:"false"`
}

type UpsParams struct {
	InputAcVoltage  float32          `json:"input_ac_voltage" example:"220"`
	InputAcCurrent  float32          `json:"input_ac_current" example:"5"`
	BatGroupVoltage float32          `json:"bat_group_voltage" example:"48"`
	BatGroupCurrent float32          `json:"bat_group_current" example:"20"`
	BatteryCapacity float32          `json:"battery_capacity" example:"20"`
	SOC             float32          `json:"soc" example:"100"` // state of charge (percent)
	Batteries       [4]BatteryParams `json:"batteries"`
	Alarms Alarms
}

func DefaultUpsParams() UpsParams {
	return UpsParams{
		InputAcVoltage:  220,
		InputAcCurrent:  5,
		BatGroupVoltage: 54,
		BatGroupCurrent: 0,
		BatteryCapacity: 50,
		SOC:             100,
		Batteries: [4]BatteryParams{
			{
				Voltage: 13.5,
				Temp:    24,
				Resist:  5,
			},
			{
				Voltage: 13.5,
				Temp:    24,
				Resist:  5,
			},
			{
				Voltage: 13.5,
				Temp:    24,
				Resist:  5,
			},
			{
				Voltage: 13.5,
				Temp:    24,
				Resist:  5,
			},
		},
	}
}
