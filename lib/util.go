package lib

import (
	"math"
	"math/rand"
)

const Grid = 12.9

// RandomRange return random value of range
func RandomRange(x, y float64) float64 {
	if x > y {
		x, y = y, x
	}

	return rand.Float64()*(y-x) + x
}

func floor(value float64, c int) float64 {
	return math.Floor(value*math.Pow10(c)) / math.Pow10(c)
}

func ccw(p1, p2, p3 Pos) int {
	temp := (p1.X*p2.Y + p2.X*p3.Y + p3.X*p1.Y) - (p1.Y*p2.X - p2.Y*p3.X - p3.Y*p1.X)
	if temp > 0 {
		return 1
	} else if temp < 0 {
		return -1
	} else {
		return 0
	}
}
