package utils_test

import (
	"math"
	"math/rand"
	"testing"

	"github.com/alex11prog/ups-imitator/internal/app/utils"
	"github.com/stretchr/testify/assert"
)

func Test_Bool2byte(t *testing.T) {
	assert.Equal(t, utils.Bool2byte(true), byte(1))
	assert.Equal(t, utils.Bool2byte(false), byte(0))
}

func Test_SimulateMeasErr(t *testing.T) {
	for i := 0; i < 100; i++ {
		src := rand.Float32() * 100
		res := utils.SimulateMeasErr(0.1, src)

		assert.GreaterOrEqual(t, 0.1, float32(math.Abs(float64((src-res)/src))))
	}
}
