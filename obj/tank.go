package obj

import (
	"math"

	"app/lib"
)

func TankTick(obj *Object) {
	obj.R = 12.9 * math.Pow(1.01, float64(obj.Level)-1)
	obj.Damage = 20 + float64(obj.Stats[2])*4
	obj.Speed = (0.07 + 0.007*float64(obj.Stats[7])) * math.Pow(0.985, float64(obj.Level)-1)

	if obj.Mh != 48.+float64(obj.Level)*2.+float64(obj.Stats[1])*20 {
		im := obj.Mh
		obj.Mh = 48. + float64(obj.Level)*2. + float64(obj.Stats[1])*20
		obj.H *= obj.Mh / im
	}

	if obj.Controller != nil {
		obj.GetDH += obj.Mh / 60 / 30 * (0.03 + 0.12*float64(obj.Stats[0]))

		for i := 0; i < 8; i++ {
			if obj.Stats[i] > obj.MaxStats[i] {
				obj.Controller.Stat += obj.Stats[i] - obj.MaxStats[i]
				obj.Stats[i] = obj.MaxStats[i]
			}
		}
	} else {
		obj.GetDH -= obj.Mh / 60 / 10
	}
}

func (o *Object) ChangeTank(c *Object) {
	o.Type = c.Type
	o.Bound = c.Bound
	o.Stance = c.Stance
	o.Sight = c.Sight
	o.MaxStats = c.MaxStats
	o.Guns = c.Guns
	o.Tick = c.Tick
	o.Collision = c.Collision
	o.KillEvent = c.KillEvent
	o.DeadEvent = c.DeadEvent
}

func NewBasic() *Object {
	var obj *Object = NewObject(map[string]interface{}{
		"type":       "Basic",
		"x":          lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
		"y":          lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
		"maxStats":   [8]int{7, 7, 7, 7, 7, 7, 7, 7},
		"isShowName": true,
	}, nil, DefaultCollision, DefaultKillEvent, nil)
	obj.Guns = []Gun{*NewGun(obj, map[string]interface{}{}, nil, nil, nil, nil)}
	return obj
}

func NewOverload() *Object {
	var obj *Object = NewObject(map[string]interface{}{
		"type":       "Overload",
		"x":          lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
		"y":          lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
		"level":      45,
		"exp":        23536,
		"stats":      [8]int{0, 0, 0, 7, 0, 7, 5, 7},
		"maxStats":   [8]int{7, 7, 7, 7, 7, 7, 7, 7},
		"sight":      1.11,
		"isShowName": true,
	}, TankTick, DefaultCollision, nil, nil)
	var dir float64 = -math.Pi / 2
	obj.Guns = []Gun{}
	for i := 0; i < 4; i++ {
		obj.Guns = append(obj.Guns, *NewGun(obj, map[string]interface{}{
			"type":     "Drone",
			"speed":    0.6,
			"damage":   0.7,
			"health":   2,
			"bound":    0.4,
			"reload":   0.167,
			"autoshot": true,
			"sy":       1.6,
			"dir":      dir,
			"limit":    2,
		}, nil, nil, nil, nil))
		dir += math.Pi / 2
	}
	return obj
}

func NewNecro() *Object {
	var obj *Object = NewObject(map[string]interface{}{
		"type":       "Necromanser",
		"x":          lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
		"y":          lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
		"level":      45,
		"exp":        23536,
		"stats":      [8]int{0, 0, 0, 7, 7, 7, 5, 7},
		"maxStats":   [8]int{7, 7, 7, 7, 7, 7, 7, 7},
		"sight":      1.11,
		"isShowName": true,
	}, TankTick, DefaultCollision, NecroKillEvent, nil)
	obj.Guns = []Gun{*NewGun(nil, map[string]interface{}{
		"limit": 0,
	}, nil, nil, nil, nil)}
	return obj
}

func NecroKillEvent(a *Object, b *Object) {
	if b.Type == "Square" && a.Guns[0].Limit < 22+2*a.Stats[6] {
		a.Guns[0].Limit++
		b.Type = "NecroSquare"
		b.Team = a.Team
		b.Exp = 0
		b.Owner = a
		b.Controller = a.Controller
		b.Mh = (8 + 6*float64(a.Stats[4])) * 0.5
		b.H = b.Mh
		b.Damage = (7 + 3*float64(a.Stats[5])) * 1.68
		b.Speed = (0.056 + 0.02*float64(a.Stats[3])) * 0.5
		b.IsDead = false
		b.IsShowHealth = false
		b.KillEvent = NecroDroneKillEvent
		b.DeadEvent = func(obj *Object, killer *Object) {
			obj.Owner.Guns[0].Limit--
			ShapeCount++
		}
		b.Tick = DefaultDroneTick
	} else {
		DefaultKillEvent(a, b)
	}
}

func NecroDroneKillEvent(a *Object, b *Object) {
	NecroKillEvent(a.Owner, b)
}
