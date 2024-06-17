package imitator

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/alex11prog/ups-imitator/internal/app/imitator/ups"
	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/goburrow/modbus"
)

type Imitator struct {
	conf          *model.Config
	upsSyncTicker *time.Ticker
	mode          atomic.Bool // true - auto, false manual
	upsSlice      []*ups.Ups
}

func New(clients []modbus.Client, conf *model.Config) *Imitator {
	res := &Imitator{
		conf:          conf,
		upsSyncTicker: time.NewTicker(conf.UpsSyncInterval),
	}
	res.mode.Store(true)
	res.upsSlice = make([]*ups.Ups, len(clients))
	for i, client := range clients {
		res.upsSlice[i] = ups.New(client, conf)
	}
	return res
}

func (im *Imitator) Start() {
	go func() {
		for range im.upsSyncTicker.C {
			for _, ups := range im.upsSlice {
				ups.CountAndSend()
			}
		}
	}()
}

func (im *Imitator) GetMode() bool {
	return im.mode.Load()
}

func (im *Imitator) SetMode(val bool) {
	old := im.mode.Swap(val)
	if old != val {
		if val {
			for _, ups := range im.upsSlice {
				ups.Reset()
			}
			im.upsSyncTicker.Reset(im.conf.UpsSyncInterval)
		} else {
			im.upsSyncTicker.Stop()
		}
	}
}

func (im *Imitator) GetUpsParams() []*model.UpsParams {
	res := make([]*model.UpsParams, len(im.upsSlice))
	for i, ups := range im.upsSlice {
		res[i] = ups.GetParams()
	}
	return res
}

func (im *Imitator) UpdateUps(ups_id int, form *model.UpsParamsUpdateForm) error {
	if l := len(im.upsSlice); ups_id >= l {
		return fmt.Errorf("ups_id out of range: %d, expected less %d", ups_id, l)
	}

	im.upsSlice[ups_id].UpdateParams(form)
	return nil
}

// UpdateUpsBattery updates ups battery
// ups_id - ups index in slice, bat_id - battery index in battery array
func (im *Imitator) UpdateUpsBattery(ups_id int, bat_id int, form *model.BatteryParamsUpdateForm) error {
	if l := len(im.upsSlice); ups_id >= l {
		return fmt.Errorf("ups_id out of range: %d, expected less %d", ups_id, l)
	}
	return im.upsSlice[ups_id].UpdateBatteryParams(bat_id, form)
}
