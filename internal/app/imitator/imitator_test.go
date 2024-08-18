package imitator

import (
	"testing"

	"github.com/alex11prog/ups-imitator/internal/app/imitator/mockmodbus"
	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/stretchr/testify/require"
)

func Test_recalcAndSendParams(t *testing.T) {
	mockModbus := mockmodbus.New()
	imitator := New(mockModbus, model.TestConfig(t))

	imitator.recalcAndSendParams()
	sentAlarmsData := mockModbus.GetWriteMultipleCoilsQueries()
	require.Equal(t, 1, len(sentAlarmsData))
	sentAlarms := sentAlarmsData[0]
	require.Equal(t, uint16(0), sentAlarms.Address)
	require.Equal(t, uint16(3), sentAlarms.Quantity)
	require.Equal(t, 1, len(sentAlarms.Value))

	sentParamsData := mockModbus.GetWriteMultipleRegistersQueries()
	require.Equal(t, 1, len(sentParamsData))
	sentParams := sentParamsData[0]
	require.Equal(t, uint16(0), sentParams.Address)
	require.Equal(t, uint16(70), sentParams.Quantity)
	require.Equal(t, 140, len(sentParams.Value))
}
