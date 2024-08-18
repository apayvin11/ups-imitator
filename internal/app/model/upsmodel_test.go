package model_test

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/alex11prog/ups-imitator/internal/app/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BatteryParams_Update(t *testing.T) {
	batteryParams := model.BatteryParams{Voltage: 12, Temp: 22, Resist: 5.5}
	testCases := []struct {
		name       string
		updateForm model.BatteryParamsUpdateForm
		expected   model.BatteryParams
	}{
		{
			name:       "Voltage",
			updateForm: model.BatteryParamsUpdateForm{Voltage: utils.NewP(float32(13))},
			expected:   model.BatteryParams{Voltage: 13, Temp: 22, Resist: 5.5},
		},
		{
			name:       "Temp",
			updateForm: model.BatteryParamsUpdateForm{Temp: utils.NewP(float32(20))},
			expected:   model.BatteryParams{Voltage: 13, Temp: 20, Resist: 5.5},
		},
		{
			name:       "Resist",
			updateForm: model.BatteryParamsUpdateForm{Resist: utils.NewP(float32(4))},
			expected:   model.BatteryParams{Voltage: 13, Temp: 20, Resist: 4},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			batteryParams.Update(tc.updateForm)
			assert.Equal(t, tc.expected, batteryParams)
		})
	}
}

func Test_Alarms_Update(t *testing.T) {
	alarms := model.Alarms{}
	testCases := []struct {
		name       string
		updateForm model.AlarmsUpdateForm
		expected   model.Alarms
	}{
		{
			name:       "UpcInBatteryMode",
			updateForm: model.AlarmsUpdateForm{UpcInBatteryMode: utils.NewP(true)},
			expected:   model.Alarms{UpcInBatteryMode: true},
		},
		{
			name:       "LowBattery",
			updateForm: model.AlarmsUpdateForm{LowBattery: utils.NewP(true)},
			expected:   model.Alarms{UpcInBatteryMode: true, LowBattery: true},
		},
		{
			name:       "Overload",
			updateForm: model.AlarmsUpdateForm{Overload: utils.NewP(true)},
			expected:   model.Alarms{UpcInBatteryMode: true, LowBattery: true, Overload: true},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			alarms.Update(tc.updateForm)
			assert.Equal(t, tc.expected, alarms)
		})
	}
}

func Test_UpsParams_Update(t *testing.T) {
	testCases := []struct {
		name       string
		src        func() *model.UpsParams
		updateForm model.UpsParamsUpdateForm
		expected   func() *model.UpsParams
	}{
		{
			name: "InputAcVoltage",
			src: func() *model.UpsParams {
				return model.TestUpsParams(t)
			},
			updateForm: model.UpsParamsUpdateForm{
				InputAcVoltage: utils.NewP(float32(230)),
			},
			expected: func() *model.UpsParams {
				params := model.TestUpsParams(t)
				params.InputAcVoltage = 230
				return params
			},
		},
		{
			name: "InputAcCurrent",
			src: func() *model.UpsParams {
				return model.TestUpsParams(t)
			},
			updateForm: model.UpsParamsUpdateForm{
				InputAcCurrent: utils.NewP(float32(6)),
			},
			expected: func() *model.UpsParams {
				params := model.TestUpsParams(t)
				params.InputAcCurrent = 6
				return params
			},
		},
		{
			name: "BatGroupVoltage",
			src: func() *model.UpsParams {
				return model.TestUpsParams(t)
			},
			updateForm: model.UpsParamsUpdateForm{
				BatGroupVoltage: utils.NewP(float32(56)),
			},
			expected: func() *model.UpsParams {
				params := model.TestUpsParams(t)
				params.BatGroupVoltage = 56
				return params
			},
		},
		{
			name: "BatGroupCurrent",
			src: func() *model.UpsParams {
				return model.TestUpsParams(t)
			},
			updateForm: model.UpsParamsUpdateForm{
				BatGroupCurrent: utils.NewP(float32(11)),
			},
			expected: func() *model.UpsParams {
				params := model.TestUpsParams(t)
				params.BatGroupCurrent = 11
				return params
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src := tc.src()
			src.Update(tc.updateForm)
			assert.Equal(t, *tc.expected(), *src)
		})
	}
}

func Test_UpsParams_GetParamBytes(t *testing.T) {
	upsParams := model.TestUpsParams(t)
	paramBytes := upsParams.GetParamBytes()
	require.Equal(t, 140, len(paramBytes))
	receivedUpsParams := model.UpsParams{
		InputAcVoltage:  math.Float32frombits(binary.BigEndian.Uint32(paramBytes[:4])),
		InputAcCurrent:  math.Float32frombits(binary.BigEndian.Uint32(paramBytes[4:8])),
		BatGroupVoltage: math.Float32frombits(binary.BigEndian.Uint32(paramBytes[8:12])),
		BatGroupCurrent: math.Float32frombits(binary.BigEndian.Uint32(paramBytes[12:16])),
	}
	assert.Equal(t, upsParams.InputAcVoltage, receivedUpsParams.InputAcVoltage)
	assert.Equal(t, upsParams.InputAcCurrent, receivedUpsParams.InputAcCurrent)
	assert.Equal(t, upsParams.BatGroupVoltage, receivedUpsParams.BatGroupVoltage)
	assert.Equal(t, upsParams.BatGroupCurrent, receivedUpsParams.BatGroupCurrent)
	for i := range 4 {
		start := 32 * (i + 1)
		receivedUpsParams.Batteries[i].Voltage = math.Float32frombits(binary.BigEndian.Uint32(paramBytes[start:]))
		receivedUpsParams.Batteries[i].Temp = math.Float32frombits(binary.BigEndian.Uint32(paramBytes[start+4:]))
		receivedUpsParams.Batteries[i].Resist = math.Float32frombits(binary.BigEndian.Uint32(paramBytes[start+8:]))
		assert.Equal(t, upsParams.Batteries[i], receivedUpsParams.Batteries[i])
	}
}

func Test_UpsParams_GetAlarmBytes(t *testing.T) {
	upsParams := model.UpsParams{
		Alarms: model.Alarms{
			UpcInBatteryMode: true,
			LowBattery:       false,
			Overload:         true,
		},
	}
	assert.Equal(t, []byte{0b00000101}, upsParams.GetAlarmBytes())
}
