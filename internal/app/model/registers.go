package model

const (
	// Holding Registers:
	RegInputAcVoltage      uint16 = 0x0000
	RegInputAcCurrent      uint16 = 0x0002
	RegBatteryGroupVoltage uint16 = 0x0004
	RegBatteryGroupCurrent uint16 = 0x0006
	RegBattery1Voltage     uint16 = 0x0010
	RegBattery1Temp        uint16 = 0x0012
	RegBattery1Res         uint16 = 0x0014
	RegBattery2Voltage     uint16 = 0x0020
	RegBattery2Temp        uint16 = 0x0022
	RegBattery2Res         uint16 = 0x0024
	RegBattery3Voltage     uint16 = 0x0030
	RegBattery3Temp        uint16 = 0x0032
	RegBattery3Res         uint16 = 0x0034
	RegBattery4Voltage     uint16 = 0x0040
	RegBattery4Temp        uint16 = 0x0042
	RegBattery4Res         uint16 = 0x0044

	// Coils
	// Alarms
	RegAlarmUpcInBatteryMode = 0x0000
	RegAlarmLowBattery       = 0x0001
	RegAlarmOverload         = 0x0002
	NumOfAlarm               = 3
)
