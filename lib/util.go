package lib

import (
	"math/rand"
)

// RandomRange return random value of range
func RandomRange(x, y float64) float64 {
	if x > y {
		x, y = y, x
	}

	return rand.Float64()*(y-x) + x
}
