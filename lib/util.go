package lib

import (
	"math"
	"math/rand"
	"time"
)

const Grid = 12.9

// Pos is
type Pos struct {
	X float64
	Y float64
}

// Score is
type Score struct {
	Name  string
	Type  string
	Score int
}

// Scoreboard is
type Scoreboard [11]Score

// Scoreboard Push
func (sc *Scoreboard) Push(value Score) {
	var index = 0
	for ; value.Score < sc[index+1].Score && index < 10; index++ {
	}
	for j := 9; j >= index; j-- {
		sc[j+1] = sc[j]
	}
	sc[index] = value
}

// RandomRange return random value of range
func RandomRange(x, y float64) float64 {
	if x > y {
		x, y = y, x
	}

	return rand.Float64()*(y-x) + x
}

func Floor(value float64, c int) float64 {
	return math.Floor(value*math.Pow10(c)) / math.Pow10(c)
}

func Lerp(a, b float64, per float64) float64 {
	return (1-per)*a + per*b
}

func Ccw(p1, p2, p3 Pos) int {
	temp := (p1.X*p2.Y + p2.X*p3.Y + p3.X*p1.Y) - (p1.Y*p2.X - p2.Y*p3.X - p3.Y*p1.X)
	if temp > 0 {
		return 1
	} else if temp < 0 {
		return -1
	} else {
		return 0
	}
}

func Now() int64 {
	return time.Now().UnixNano() / 1000000
}
