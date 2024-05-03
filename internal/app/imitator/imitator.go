package imitator

import (
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
	ups_          []*ups.Ups
}

func New(clients []modbus.Client, conf *model.Config) *Imitator {
	res := &Imitator{
		conf:          conf,
		upsSyncTicker: time.NewTicker(conf.UpsSyncInterval),
	}
	res.mode.Store(true)
	for _, client := range clients {
		res.ups_ = append(res.ups_, ups.New(client, conf))
	}
	return res
}

func (im *Imitator) Start() {
	for range im.upsSyncTicker.C {
		for _, ups := range im.ups_ {
			ups.CountAndSend()
		}
	}
}

func (im *Imitator) GetMode() bool {
	return im.mode.Load()
}

func (im *Imitator) SetMode(val bool) {
	old := im.mode.Swap(val)
	if old != val {
		if val {
			for _, ups := range im.ups_ {
				ups.Reset()
			}
			im.upsSyncTicker.Reset(im.conf.UpsSyncInterval)
		} else {
			im.upsSyncTicker.Stop()
		}
	}
}

func (im *Imitator) GetUpsParams() []*model.UpsParams {
	var res []*model.UpsParams
	for _, ups := range im.ups_ {
		res = append(res, ups.GetParams())
	}
	return res
}
