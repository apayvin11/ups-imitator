package utils

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SimulateMeasErr(t *testing.T) {
	for i := 0; i < 100; i++ {
		src := rand.Float32() * 100
		res := SimulateMeasErr(0.1, src)

		assert.GreaterOrEqual(t, 0.1, float32(math.Abs(float64((src-res)/src))))
	}
}
