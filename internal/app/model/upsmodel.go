package model

import (
	"encoding/binary"
	"math"

	"github.com/alex11prog/ups-imitator/internal/app/utils"
)

type BatteryParams struct {
	Voltage float32 `json:"voltage" example:"12"`
	Temp    float32 `json:"temp" example:"24"`
	Resist  float32 `json:"resist" example:"5"`
}

func (bat *BatteryParams) Update(form BatteryParamsUpdateForm) {
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

func (a *Alarms) Update(form AlarmsUpdateForm) {
	if form.UpcInBatteryMode != nil {
		a.UpcInBatteryMode = *form.UpcInBatteryMode
	}
	if form.LowBattery != nil {
		a.LowBattery = *form.LowBattery
	}
	if form.Overload != nil {
		a.Overload = *form.Overload
	}
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

	Alarms Alarms `json:"alarms"`
}

func (ups *UpsParams) Update(form UpsParamsUpdateForm) {
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
}

type UpsParamsUpdateForm struct {
	InputAcVoltage  *float32 `json:"input_ac_voltage" example:"220"` // V
	InputAcCurrent  *float32 `json:"input_ac_current" example:"5"`   // Amp
	BatGroupVoltage *float32 `json:"bat_group_voltage" example:"48"` // V
	BatGroupCurrent *float32 `json:"bat_group_current" example:"0"`  // Amp
}

func (ups *UpsParams) GetParamBytes() []byte {
	res := make([]byte, RegBattery4Res*2 +4 )
	binary.BigEndian.PutUint32(res[RegInputAcVoltage*2:], math.Float32bits(ups.InputAcVoltage))
	binary.BigEndian.PutUint32(res[RegInputAcCurrent*2:], math.Float32bits(ups.InputAcCurrent))
	binary.BigEndian.PutUint32(res[RegBatteryGroupVoltage*2:], math.Float32bits(ups.BatGroupVoltage))
	binary.BigEndian.PutUint32(res[RegBatteryGroupCurrent*2:], math.Float32bits(ups.BatGroupCurrent))
	for i, battery := range ups.Batteries {
		start := RegBattery1Voltage * uint16(2) * (uint16(i) + 1)
		binary.BigEndian.PutUint32(res[start:], math.Float32bits(battery.Voltage))
		binary.BigEndian.PutUint32(res[start+4:], math.Float32bits(battery.Temp))
		binary.BigEndian.PutUint32(res[start+8:], math.Float32bits(battery.Resist))
	}
	return res
}

func (ups *UpsParams) GetAlarmBytes() []byte {
	var alarmsBitwise byte
	alarmsBitwise |= utils.Bool2byte(ups.Alarms.UpcInBatteryMode)
	alarmsBitwise |= utils.Bool2byte(ups.Alarms.LowBattery) << 1
	alarmsBitwise |= utils.Bool2byte(ups.Alarms.Overload) << 2
	return []byte{alarmsBitwise}
}
