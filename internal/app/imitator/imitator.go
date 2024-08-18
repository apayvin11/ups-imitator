package imitator

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/alex11prog/ups-imitator/internal/app/imitator/ups"
	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/goburrow/modbus"
)

type Imitator struct {
	client        modbus.Client
	conf          *model.Config
	upsSyncTicker *time.Ticker
	mode          atomic.Bool // true - auto, false - manual
	ups           *ups.Ups
}

func New(client modbus.Client, conf *model.Config) *Imitator {
	res := &Imitator{
		client:        client,
		conf:          conf,
		upsSyncTicker: time.NewTicker(conf.UpsSyncInterval),
		ups:           ups.New(conf),
	}
	res.mode.Store(true)
	return res
}

// Start starts working in the background, recalculating and sending parameters to the UPS via Modbus
func (im *Imitator) Start() {
	go func() {
		for range im.upsSyncTicker.C {
			im.recalcAndSendParams()
		}
	}()
}

func (im *Imitator) recalcAndSendParams() {
	im.ups.RecalculateParams()
	params := im.ups.GetParamsWithSimulatedMeasErr()
	paramBytes := params.GetParamBytes()
	if _, err := im.client.WriteMultipleRegisters(model.RegInputAcVoltage, uint16(len(paramBytes)/2), paramBytes); err != nil {
		log.Println(err)
		return
	}
	log.Printf("InputAcVoltage: %v\n", params.InputAcVoltage)
	log.Printf("InputAcCurrent: %v\n", params.InputAcCurrent)
	log.Printf("BatGroupVoltage: %v\n", params.BatGroupVoltage)
	log.Printf("BatGroupCurrent: %v\n", params.BatGroupCurrent)
	log.Printf("LoadCurrent: %v\n", params.LoadCurrent)
	log.Printf("RemainingBatCapacity: %v\n", params.RemainingBatCapacity)
	log.Printf("SOC: %v\n\n", params.SOC)

	alarmBytes := params.GetAlarmBytes()
	if _, err := im.client.WriteMultipleCoils(model.RegAlarmUpcInBatteryMode, model.NumOfAlarm, alarmBytes); err != nil {
		log.Println(err)
	}
}

func (im *Imitator) GetMode() bool {
	return im.mode.Load()
}

func (im *Imitator) SetMode(val bool) {
	old := im.mode.Swap(val)
	if old != val {
		if val {
			im.ups.Reset()
			im.upsSyncTicker.Reset(im.conf.UpsSyncInterval)
		} else {
			im.upsSyncTicker.Stop()
		}
	}
}

func (im *Imitator) GetAllUpsParams() model.UpsParams {
	return im.ups.GetAllParams()
}

func (im *Imitator) UpdateUpsParams(form model.UpsParamsUpdateForm) {
	im.ups.UpdateParams(form)
}

func (im *Imitator) UpdateUpsBatteryParams(bat_id int, batParams model.BatteryParamsUpdateForm) error {
	return im.ups.UpdateBatteryParams(bat_id, batParams)
}

func (im *Imitator) UpdateAlarms(alarms model.AlarmsUpdateForm) {
	im.ups.UpdateAlarms(alarms)
}
