package ups

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/alex11prog/ups-imitator/internal/app/utils"
)

type chargeState uint8

const (
	chargedState chargeState = iota
	dischargingState
	dischargedState
	chargingState
)

type Ups struct {
	conf *model.Config

	mu             sync.Mutex
	state          chargeState
	lastUpdateTime time.Time
	cycleDoneTime  time.Time // charge or discharge
	params         model.UpsParams
}

func New(conf *model.Config) *Ups {
	u := &Ups{
		conf:           conf,
		lastUpdateTime: time.Now(),
		cycleDoneTime:  time.Now(),
	}
	u.setDefaultUpsParams()
	return u
}

// Reset sets the default value for all params
func (u *Ups) Reset() {
	u.mu.Lock()
	u.state = chargedState
	u.lastUpdateTime = time.Now()
	u.cycleDoneTime = time.Now()
	u.setDefaultUpsParams()
	u.mu.Unlock()
}

// RecalculateParams recalculates parameters depending on the ups state
func (u *Ups) RecalculateParams() {
	u.mu.Lock()
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
	u.mu.Unlock()
}

func (u *Ups) GetAllParams() (params model.UpsParams) {
	u.mu.Lock()
	params = u.params
	u.mu.Unlock()
	return
}

func (u *Ups) GetParamsWithSimulatedMeasErr() (params model.UpsParams) {
	u.mu.Lock()
	params = model.UpsParams{
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
		params.Batteries[i].Voltage = utils.SimulateMeasErr(0.04, bat.Voltage)
		params.Batteries[i].Temp = utils.SimulateMeasErr(0.04, bat.Temp)
		params.Batteries[i].Resist = utils.SimulateMeasErr(0.04, bat.Resist)
	}
	u.mu.Unlock()
	return
}

func (u *Ups) UpdateParams(params model.UpsParamsUpdateForm) {
	u.mu.Lock()
	u.params.Update(params)
	u.mu.Unlock()
}

func (u *Ups) UpdateBatteryParams(bat_id int, batParams model.BatteryParamsUpdateForm) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	if l := len(u.params.Batteries); bat_id >= l {
		return fmt.Errorf("bat_id out of range: %d, expected less %d", bat_id, l)
	}
	u.params.Batteries[bat_id].Update(batParams)
	return nil
}

func (u *Ups) UpdateAlarms(alarms model.AlarmsUpdateForm) {
	u.mu.Lock()
	u.params.Alarms.Update(alarms)
	u.mu.Unlock()
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

func (u *Ups) setState(s chargeState) {
	log.Printf("\nnew state: %v\n\n", s)
	u.state = s
}

func (u *Ups) recalcLoadCurrent() {
	u.params.LoadCurrent = u.conf.LoadPower / u.params.BatGroupVoltage
}

func (u *Ups) recalcSoc() {
	u.params.SOC = u.params.RemainingBatCapacity / u.params.BatCapacity
}

// recalcBatGroupVoltage recalculates BatGroupVoltage depending on battery current and SOC (state of charge)
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
