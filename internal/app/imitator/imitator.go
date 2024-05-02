package imitator

import (
	"sync/atomic"
	"time"

	"github.com/alex11prog/ups-imitator/internal/app/imitator/ups"
	"github.com/goburrow/modbus"
)

type Imitator struct {
	upsSyncInterval time.Duration
	mode            atomic.Bool // true - auto, false manual
	ups_            []*ups.Ups
}

func New(clients []modbus.Client, upsSyncInterval time.Duration) *Imitator {
	res := &Imitator{
		upsSyncInterval: upsSyncInterval,
	}
	res.mode.Store(true)
	for _, client := range clients {
		res.ups_ = append(res.ups_, ups.New(client))
	}
	return res
}

func (im *Imitator) Start() {
	ticker := time.NewTicker(im.upsSyncInterval)
	for range ticker.C {
		if !im.mode.Load() {
			continue
		}
		for _, ups := range im.ups_ {
			ups.CountAndSend()
		}
	}
}

func (im *Imitator) GetMode() bool {
	return im.mode.Load()
}

func (im *Imitator) SetMode(val bool) {
	im.mode.Store(val)
}
