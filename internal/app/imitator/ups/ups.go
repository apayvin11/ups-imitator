package ups

import (
	"sync"

	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/goburrow/modbus"
)

type Ups struct {
	client modbus.Client

	mu     sync.Mutex
	params model.UpsParams
}

func New(client modbus.Client) *Ups {
	return &Ups{
		client: client,
		params: model.DefaultUpsParams(),
	}
}

// CountAndSend recalculates and sends new values ​​via modbus to UPS
func (u *Ups) CountAndSend() error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// holding registers
	/*{
		var buf bytes.Buffer
		if err := binary.Write(&buf, binary.BigEndian, u.params.InputAcVoltage); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}
		if err := binary.Write(&buf, binary.BigEndian, u.params.InputAcCurrent); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}
		if err := binary.Write(&buf, binary.BigEndian, u.params.BatGroupVoltage); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}
		if err := binary.Write(&buf, binary.BigEndian, u.params.BatGroupCurrent); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}
		res := buf.Bytes()
		if _, err := u.client.WriteMultipleRegisters(regInputAcVoltage, uint16(len(res)/2), res); err != nil {
			return err
		}
	}*/
	// coils
	/* 	{
		var buf bytes.Buffer
		if err := binary.Write(&buf, binary.BigEndian, u.params.Alarms.UpcInBatteryMode); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}
		if err := binary.Write(&buf, binary.BigEndian, u.params.Alarms.LowBattery); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}
		if err := binary.Write(&buf, binary.BigEndian, u.params.Alarms.Overload); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}
		res := buf.Bytes()
		if _, err := u.client.WriteMultipleCoils(regAlarmUpcInBatteryMode, uint16(len(res)), res); err != nil {
			return err
		}
	} */
	res := byte(0b11100101)
	if _, err := u.client.WriteMultipleCoils(regAlarmUpcInBatteryMode, 3, []byte{res}); err != nil {
		return err
	}
	return nil
}
