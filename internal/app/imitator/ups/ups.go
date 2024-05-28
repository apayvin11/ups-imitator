package ups

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/alex11prog/ups-imitator/internal/app/utils"
	"github.com/goburrow/modbus"
)

type chargeState uint8

const (
	chargedState chargeState = iota
	dischargingState
	dischargedState
	chargingState
)

type Ups struct {
	client modbus.Client
	conf   *model.Config

	mu             sync.Mutex
	state          chargeState
	lastUpdateTime time.Time
	cycleDoneTime  time.Time // charge or discharge
	params         model.UpsParams
}

func New(client modbus.Client, conf *model.Config) *Ups {
	u := &Ups{
		client:         client,
		conf:           conf,
		lastUpdateTime: time.Now(),
		cycleDoneTime:  time.Now(),
	}
	u.setDefaultUpsParams()
	return u
}

func (u *Ups) Reset() {
	u.mu.Lock()
	u.state = chargedState
	u.lastUpdateTime = time.Now()
	u.cycleDoneTime = time.Now()
	u.setDefaultUpsParams()
	u.mu.Unlock()
}

// CountAndSend recalculates and sends new values ​​via modbus to UPS
func (u *Ups) CountAndSend() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.recalculateParams()
	return u.sendParams(u.getParamsWithSimulatedMeasErr())
}

func (u *Ups) setState(s chargeState) {
	log.Printf("\nnew state: %v\n\n", s)
	u.state = s
}

func (u *Ups) recalculateParams() {
	switch u.state {
	case chargedState:
		if time.Since(u.cycleDoneTime) > u.conf.CycleChangeTimeout {
			u.params.InputAcVoltage = 0
			u.params.InputAcCurrent = 0
			u.recalcLoadCurrent()
			u.params.BatGroupCurrent = -u.params.LoadCurrent * 1.1
			u.params.Alarms.UpcInBatteryMode = true

			u.setState(dischargingState)
		}

	case dischargingState:
		elapsedTimeH := float32(time.Since(u.lastUpdateTime)) / float32(time.Hour) // elapsed time in hours
		spentCapacity := u.params.BatGroupCurrent * elapsedTimeH                   // Ah
		u.params.RemainingBatCapacity += spentCapacity

		if u.params.RemainingBatCapacity < 0 {
			u.params.RemainingBatCapacity = 0
			u.params.SOC = 0
			u.cycleDoneTime = time.Now()
			u.params.LoadCurrent = 0
			u.params.BatGroupCurrent = 0
			u.setState(dischargedState)
			break
		}
		u.recalcSoc()
		u.recalcBatGroupVoltage()
		u.recalcLoadCurrent()
		u.params.BatGroupCurrent = -u.params.LoadCurrent * 1.1
		u.recalcInputAcCurrent()

		if u.params.SOC < u.conf.LowSocTriggerAlarm {
			u.params.Alarms.LowBattery = true
		}

	case dischargedState:
		if time.Since(u.cycleDoneTime) > u.conf.CycleChangeTimeout {
			u.params.InputAcVoltage = u.conf.DefaultInputAcVoltage
			u.params.BatGroupCurrent = u.conf.ChargeCurrentLimit
			u.recalcLoadCurrent()
			u.recalcInputAcCurrent()
			u.params.Alarms = model.Alarms{}

			u.setState(chargingState)
		}

	case chargingState:
		elapsedTimeH := float32(time.Since(u.lastUpdateTime)) / float32(time.Hour) // elapsed time in hours
		receivedCapacity := u.params.BatGroupCurrent * elapsedTimeH                // Ah
		u.params.RemainingBatCapacity += receivedCapacity

		if u.params.RemainingBatCapacity > u.params.BatCapacity {
			u.params.RemainingBatCapacity = u.params.BatCapacity
			u.params.SOC = 1
			u.cycleDoneTime = time.Now()
			u.params.BatGroupCurrent = 0

			u.setState(chargedState)
			break
		}

		u.recalcSoc()
		u.recalcBatGroupVoltage()
		u.recalcLoadCurrent()
		u.recalcChargingCurrent()
		u.recalcInputAcCurrent()
	}
	u.recalcBatValtages()
	u.lastUpdateTime = time.Now()
}

// sendParams sends new values ​​via modbus to UPS
func (u *Ups) sendParams(params *model.UpsParams) error {
	log.Printf("InputAcVoltage: %v\n", params.InputAcVoltage)
	log.Printf("InputAcCurrent: %v\n", params.InputAcCurrent)
	log.Printf("BatGroupVoltage: %v\n", params.BatGroupVoltage)
	log.Printf("BatGroupCurrent: %v\n", params.BatGroupCurrent)
	log.Printf("LoadCurrent: %v\n", params.LoadCurrent)
	log.Printf("RemainingBatCapacity: %v\n", params.RemainingBatCapacity)
	log.Printf("SOC: %v\n\n", params.SOC)
	{ // holding registers
		var buf bytes.Buffer

		if err := binary.Write(&buf, binary.BigEndian, params.InputAcVoltage); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}

		if err := binary.Write(&buf, binary.BigEndian, params.InputAcCurrent); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}

		if err := binary.Write(&buf, binary.BigEndian, params.BatGroupVoltage); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}

		if err := binary.Write(&buf, binary.BigEndian, params.BatGroupCurrent); err != nil {
			return fmt.Errorf("binary.Write failed: %v", err)
		}

		buf.Write(make([]byte, 16)) // skip empty registers

		for _, battery := range params.Batteries {
			if err := binary.Write(&buf, binary.BigEndian, battery.Voltage); err != nil {
				return fmt.Errorf("binary.Write failed: %v", err)
			}

			if err := binary.Write(&buf, binary.BigEndian, battery.Temp); err != nil {
				return fmt.Errorf("binary.Write failed: %v", err)
			}

			if err := binary.Write(&buf, binary.BigEndian, battery.Resist); err != nil {
				return fmt.Errorf("binary.Write failed: %v", err)
			}

			buf.Write(make([]byte, 20)) // skip empty registers
		}

		res := buf.Bytes()
		if _, err := u.client.WriteMultipleRegisters(regInputAcVoltage, uint16(len(res)/2), res); err != nil {
			return err
		}
	}
	{ // coils
		var res byte
		res |= utils.Bool2byte(params.Alarms.UpcInBatteryMode)
		res |= utils.Bool2byte(params.Alarms.LowBattery) << 1
		res |= utils.Bool2byte(params.Alarms.Overload) << 2
		if _, err := u.client.WriteMultipleCoils(regAlarmUpcInBatteryMode, 3, []byte{res}); err != nil {
			return err
		}
	}
	return nil
}

// GetParams return params of all ups
func (u *Ups) GetParams() *model.UpsParams {
	u.mu.Lock()
	cp := u.params
	u.mu.Unlock()
	return &cp
}

func (u *Ups) UpdateParams(params *model.UpsParamsUpdateForm) {
	u.mu.Lock()
	u.params.UpdateParams(params)
	u.mu.Unlock()
}

func (u *Ups) UpdateBatteryParams(bat_id int, batParams *model.BatteryParamsUpdateForm) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if l := len(u.params.Batteries); bat_id >= l {
		return fmt.Errorf("bat_id out of range: %d, expected less %d", bat_id, l)
	}
	u.params.Batteries[bat_id].Update(batParams)
	return nil
}

func (u *Ups) getParamsWithSimulatedMeasErr() *model.UpsParams {
	res := &model.UpsParams{
		InputAcVoltage:       utils.SimulateMeasErr(0.02, u.params.InputAcVoltage),
		InputAcCurrent:       utils.SimulateMeasErr(0.02, u.params.InputAcCurrent),
		BatGroupVoltage:      utils.SimulateMeasErr(0.02, u.params.BatGroupVoltage),
		BatGroupCurrent:      utils.SimulateMeasErr(0.02, u.params.BatGroupCurrent),
		LoadCurrent:          utils.SimulateMeasErr(0.02, u.params.LoadCurrent),
		BatCapacity:          u.params.BatCapacity,
		RemainingBatCapacity: u.params.RemainingBatCapacity,
		SOC:                  u.params.SOC,
		Alarms:               u.params.Alarms,
	}
	for i, bat := range u.params.Batteries {
		res.Batteries[i].Voltage = utils.SimulateMeasErr(0.04, bat.Voltage)
		res.Batteries[i].Temp = utils.SimulateMeasErr(0.04, bat.Temp)
		res.Batteries[i].Resist = utils.SimulateMeasErr(0.04, bat.Resist)
	}
	return res
}

func (u *Ups) setDefaultUpsParams() {
	u.params = model.UpsParams{
		InputAcVoltage:       u.conf.DefaultInputAcVoltage,
		InputAcCurrent:       u.conf.LoadPower * 1.1 / u.conf.DefaultInputAcVoltage,
		BatGroupVoltage:      u.conf.MaxBatGroupVoltage,
		BatGroupCurrent:      0,
		LoadCurrent:          u.conf.LoadPower / u.conf.MaxBatGroupVoltage,
		BatCapacity:          u.conf.DefaultBatCapacity,
		RemainingBatCapacity: u.conf.DefaultBatCapacity,
		SOC:                  1,
		Batteries: [4]model.BatteryParams{
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

func (u *Ups) recalcLoadCurrent() {
	u.params.LoadCurrent = u.conf.LoadPower / u.params.BatGroupVoltage
}

func (u *Ups) recalcSoc() {
	u.params.SOC = u.params.RemainingBatCapacity / u.params.BatCapacity
}

func (u *Ups) recalcBatGroupVoltage() {
	if u.params.BatGroupCurrent < 0 { // discharge
		u.params.BatGroupVoltage = u.conf.MinBatGroupVoltage + u.params.SOC*(u.conf.MaxBatGroupVoltage-u.conf.MinBatGroupVoltage)
	} else { //charge
		u.params.BatGroupVoltage = u.conf.MinBatGroupVoltage + 1.25*u.params.SOC*(u.conf.MaxBatGroupVoltage-u.conf.MinBatGroupVoltage)
	}
	if u.params.BatGroupVoltage > u.conf.MaxBatGroupVoltage {
		u.params.BatGroupVoltage = u.conf.MaxBatGroupVoltage
	}
}

func (u *Ups) recalcInputAcCurrent() {
	if u.params.InputAcVoltage == 0 {
		u.params.InputAcCurrent = 0
		return
	}
	totalower := 1.1 * (u.conf.LoadPower + u.params.BatGroupVoltage*u.params.BatGroupCurrent)
	u.params.InputAcCurrent = totalower / u.params.InputAcVoltage
}

func (u *Ups) recalcChargingCurrent() {
	if u.params.SOC < 0.8 {
		u.params.BatGroupCurrent = u.conf.ChargeCurrentLimit
	} else {
		cf := 1 - 4*(u.params.SOC-0.8)
		u.params.BatGroupCurrent = u.conf.ChargeCurrentLimit * cf
	}
}

func (u *Ups) recalcBatValtages() {
	avgBatVoltage := u.params.BatGroupVoltage / float32(len(u.params.Batteries))
	for i := range u.params.Batteries {
		u.params.Batteries[i].Voltage = avgBatVoltage
	}
}
