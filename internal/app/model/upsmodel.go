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
	InputAcVoltage       float32          `json:"input_ac_voltage" example:"220"`          // V
	InputAcCurrent       float32          `json:"input_ac_current" example:"5"`            // Amp
	BatGroupVoltage      float32          `json:"bat_group_voltage" example:"48"`          // V
	BatGroupCurrent      float32          `json:"bat_group_current" example:"0"`           // Amp
	LoadCurrent          float32          `json:"load_current" example:"20"`               // Amp
	BatCapacity          float32          `json:"battery_capacity" example:"50"`           // Ah
	RemainingBatCapacity float32          `json:"remaining_battery_capacity" example:"50"` // Ah
	SOC                  float32          `json:"soc" example:"100"`                       // state of charge (percent)
	Batteries            [4]BatteryParams `json:"batteries"`
	Alarms               Alarms
}