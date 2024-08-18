package ups

import (
	"testing"

	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func Test_RecalculateParams(t *testing.T) {
	conf := model.TestConfig(t)
	ups := New(conf)
	ups.RecalculateParams()
	assert.Equal(t, chargedState, ups.state)

	ups.cycleDoneTime = ups.cycleDoneTime.Add(-conf.CycleChangeTimeout * 2)
	ups.RecalculateParams()
	assert.Equal(t, dischargingState, ups.state)

	ups.params.RemainingBatCapacity = -1
	ups.RecalculateParams()
	assert.Equal(t, dischargedState, ups.state)

	ups.cycleDoneTime = ups.cycleDoneTime.Add(-conf.CycleChangeTimeout * 2)
	ups.RecalculateParams()
	assert.Equal(t, chargingState, ups.state)

	ups.params.RemainingBatCapacity = conf.DefaultBatCapacity + 1
	ups.RecalculateParams()
	assert.Equal(t, chargedState, ups.state)
}
