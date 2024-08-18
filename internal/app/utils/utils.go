package utils

import "math/rand"

func Bool2byte(val bool) byte {
	if val {
		return 1
	}
	return 0
}

// SimulateMeasErr simulates measure error
// errDev - err deviation. for example 0.1 (percent)
func SimulateMeasErr(errDev, srcVal float32) float32 {
	measErrCoef := rand.Float32() * 2 - 1 // from - 1 to 1
	return srcVal + measErrCoef * (srcVal * errDev)
}

// NewP returns pointer of value
func NewP[V any](v V) *V {
	return &v
}