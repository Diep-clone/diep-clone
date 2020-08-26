package obj

import (
	"math"

	"app/lib"
)

var ShapeCount int

// TODO : other Shape and Bosses
func AddShape() {
	for ; ShapeCount > 0; ShapeCount-- {
		Objects = append(Objects, NewObject(map[string]interface{}{
			"type":   "Square",
			"name":   "Square",
			"team":   "shape",
			"mh":     10,
			"h":      10,
			"x":      lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
			"y":      lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
			"dir":    lib.RandomRange(-math.Pi, math.Pi),
			"stance": 0.1,
			"exp":    10,
		}, DefaultShapeTick, DefaultCollision, DefaultKillEvent, func(o *Object, killer *Object) {
			ShapeCount++
		}))
	}
}

// TODO : shape AI
func DefaultShapeTick(o *Object) {
	o.Dir += 0.005
}
