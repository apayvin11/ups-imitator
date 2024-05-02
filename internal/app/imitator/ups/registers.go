package ups

const (
	// Holding Registers:
	regInputAcVoltage      uint16 = 0x0000
	regInputAcCurrent      uint16 = 0x0002
	regBatteryGroupVoltage uint16 = 0x0004
	regBatteryGroupCurrent uint16 = 0x0006
	regBattery1Voltage     uint16 = 0x0010
	regBattery1Temp        uint16 = 0x0012
	regBattery1Res         uint16 = 0x0014
	regBattery2Voltage     uint16 = 0x0020
	regBattery2Temp        uint16 = 0x0022
	regBattery2Res         uint16 = 0x0024
	regBattery3Voltage     uint16 = 0x0030
	regBattery3Temp        uint16 = 0x0032
	regBattery3Res         uint16 = 0x0034
	regBattery4Voltage     uint16 = 0x0040
	regBattery4Temp        uint16 = 0x0042
	regBattery4Res         uint16 = 0x0044

	// Coils
	// Alarms
	regAlarmUpcInBatteryMode = 0x0000
	regAlarmLowBattery       = 0x0001
	regAlarmOverload         = 0x0002
)
