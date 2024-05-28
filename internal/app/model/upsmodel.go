package model

type BatteryParams struct {
	Voltage float32 `json:"voltage" example:"12"`
	Temp    float32 `json:"temp" example:"24"`
	Resist  float32 `json:"resist" example:"5"`
}

func (bat *BatteryParams) Update(form *BatteryParamsUpdateForm) {
	if form.Voltage != nil {
		bat.Voltage = *form.Voltage
	}
	if form.Temp != nil {
		bat.Temp = *form.Temp
	}
	if form.Resist != nil {
		bat.Resist = *form.Resist
	}
}

type BatteryParamsUpdateForm struct {
	Voltage *float32 `json:"voltage" example:"12"`
	Temp    *float32 `json:"temp" example:"24"`
	Resist  *float32 `json:"resist" example:"5"`
}

type Alarms struct {
	UpcInBatteryMode bool `json:"upc_in_battery_mode" example:"false"`
	LowBattery       bool `json:"low_battery" example:"false"`
	Overload         bool `json:"overload" example:"false"`
}

type AlarmsUpdateForm struct {
	UpcInBatteryMode *bool `json:"upc_in_battery_mode" example:"false"`
	LowBattery       *bool `json:"low_battery" example:"false"`
	Overload         *bool `json:"overload" example:"false"`
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
	Alarms               Alarms           `json:"alarms"`
}

func (ups *UpsParams) UpdateParams(form *UpsParamsUpdateForm) {
	if form.InputAcVoltage != nil {
		ups.InputAcVoltage = *form.InputAcVoltage
	}
	if form.InputAcCurrent != nil {
		ups.InputAcCurrent = *form.InputAcCurrent
	}
	if form.BatGroupVoltage != nil {
		ups.BatGroupVoltage = *form.BatGroupVoltage
	}
	if form.BatGroupCurrent != nil {
		ups.BatGroupCurrent = *form.BatGroupCurrent
	}
	if form.Alarms.UpcInBatteryMode != nil {
		ups.Alarms.UpcInBatteryMode = *form.Alarms.UpcInBatteryMode
	}
	if form.Alarms.LowBattery != nil {
		ups.Alarms.LowBattery = *form.Alarms.LowBattery
	}
	if form.Alarms.Overload != nil {
		ups.Alarms.Overload = *form.Alarms.Overload
	}
}

type UpsParamsUpdateForm struct {
	InputAcVoltage  *float32 `json:"input_ac_voltage" example:"220"` // V
	InputAcCurrent  *float32 `json:"input_ac_current" example:"5"`   // Amp
	BatGroupVoltage *float32 `json:"bat_group_voltage" example:"48"` // V
	BatGroupCurrent *float32 `json:"bat_group_current" example:"0"`  // Amp
	Alarms          AlarmsUpdateForm
}
