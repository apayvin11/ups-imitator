package model_test

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func Test_BatteryParams_Update(t *testing.T) {

}

func Test_Alarms_Update(t *testing.T) {

}

func Test_UpsParams_Update(t *testing.T) {

}

func Test_UpsParams_GetParamBytes(t *testing.T) {
	upsParams := model.TestUpsParams(t)
	expected := make([]byte, 140)
	binary.BigEndian.PutUint32(expected[:4], math.Float32bits(upsParams.InputAcVoltage))
	binary.BigEndian.PutUint32(expected[4:8], math.Float32bits(upsParams.InputAcCurrent))
	binary.BigEndian.PutUint32(expected[8:12], math.Float32bits(upsParams.BatGroupVoltage))
	binary.BigEndian.PutUint32(expected[12:16], math.Float32bits(upsParams.BatGroupCurrent))
	for i, battery := range upsParams.Batteries {
		start := 32 * (i + 1)
		binary.BigEndian.PutUint32(expected[start:], math.Float32bits(battery.Voltage))
		binary.BigEndian.PutUint32(expected[start+4:], math.Float32bits(battery.Temp))
		binary.BigEndian.PutUint32(expected[start+8:], math.Float32bits(battery.Resist))
	}
	assert.Equal(t, expected, upsParams.GetParamBytes())
}

func Test_UpsParams_GetAlarmBytes(t *testing.T) {

}
