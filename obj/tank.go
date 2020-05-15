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

func NewTank(t string) *Object {
	var obj *Object = NewObject(map[string]interface{}{
		"type":       t,
		"x":          lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
		"y":          lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
		"level":      45,
		"exp":        23536,
		"stats":      [8]int{0, 0, 0, 5, 7, 7, 7, 7},
		"maxStats":   [8]int{7, 7, 7, 7, 7, 7, 7, 7},
		"isShowName": true,
	}, TankTick, DefaultCollision, DefaultKillEvent, nil)

	switch t {
	case "Basic":
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{}, nil, nil, nil, nil)}
	case "Twin":
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"damage":   0.65,
			"gunbound": 0.75,
			"sx":       0.5,
			"sy":       1.88,
		}, nil, nil, nil, nil), NewGun(obj, map[string]interface{}{
			"damage":   0.65,
			"gunbound": 0.75,
			"waittime": 0.5,
			"sx":       -0.5,
			"sy":       1.88,
		}, nil, nil, nil, nil)}
	case "Triplet":
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"damage":   0.6,
			"health":   0.7,
			"gunbound": 0.5,
			"waittime": 0.5,
			"sx":       0.5,
			"sy":       1.6,
		}, nil, nil, nil, nil), NewGun(obj, map[string]interface{}{
			"damage":   0.6,
			"health":   0.7,
			"bound":    0.5,
			"waittime": 0.5,
			"sx":       -0.5,
			"sy":       1.6,
		}, nil, nil, nil, nil), NewGun(obj, map[string]interface{}{
			"damage":   0.6,
			"health":   0.7,
			"gunbound": 0.5,
		}, nil, nil, nil, nil)}
	case "TripleShot":
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"damage": 0.7,
		}, nil, nil, nil, nil), NewGun(obj, map[string]interface{}{
			"damage": 0.7,
			"dir":    math.Pi / 4,
		}, nil, nil, nil, nil), NewGun(obj, map[string]interface{}{
			"damage": 0.7,
			"dir":    -math.Pi / 4,
		}, nil, nil, nil, nil)}
	case "OctoTank":
		obj.Guns = []Gun{}
		for i := -math.Pi / 2; i <= math.Pi; i += math.Pi / 2 {
			obj.Guns = append(obj.Guns, NewGun(obj, map[string]interface{}{
				"damage": 0.65,
				"dir":    i,
			}, nil, nil, nil, nil))
		}
		for i := -math.Pi / 4 * 3; i <= math.Pi; i += math.Pi / 2 {
			obj.Guns = append(obj.Guns, NewGun(obj, map[string]interface{}{
				"damage":   0.65,
				"dir":      i,
				"waittime": 0.5,
			}, nil, nil, nil, nil))
		}
	case "Sniper":
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"speed":  1.5,
			"bound":  3,
			"rdir":   math.Pi / 120,
			"reload": 0.667,
			"sy":     2.2,
		}, nil, nil, nil, nil)}
	case "MachineGun":
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"damage": 0.7,
			"reload": 2,
			"rdir":   math.Pi / 8,
			"sy":     1.6,
		}, nil, nil, nil, nil)}
	case "FlankGuard":
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{}, nil, nil, nil, nil),
			NewGun(obj, map[string]interface{}{
				"dir": math.Pi,
				"sy":  1.6,
			}, nil, nil, nil, nil)}
	case "Destroyer":
		obj.Stats = [8]int{0, 0, 0, 7, 7, 7, 5, 7}
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"speed":    0.75,
			"damage":   3,
			"health":   2,
			"radius":   1.75,
			"gunbound": 15,
			"bound":    0.1,
			"reload":   0.25,
			"sy":       1.88,
		}, nil, nil, nil, nil)}
	case "Overlord":
		obj.Stats = [8]int{0, 0, 0, 7, 7, 7, 5, 7}
		obj.Sight = 1.11
		var dir float64 = -math.Pi / 2
		obj.Guns = []Gun{}
		for i := 0; i < 4; i++ {
			obj.Guns = append(obj.Guns, NewGun(obj, map[string]interface{}{
				"type":     "Drone",
				"speed":    0.6,
				"damage":   0.7,
				"health":   2,
				"gunbound": 0.4,
				"reload":   0.167,
				"autoshot": true,
				"sy":       1.6,
				"dir":      dir,
				"limit":    2,
			}, nil, nil, nil, nil))
			dir += math.Pi / 2
		}
	case "Necromancer":
		obj.Stats = [8]int{0, 0, 0, 7, 7, 7, 5, 7}
		obj.Sight = 1.11
		obj.KillEvent = NecroKillEvent
		obj.Guns = []Gun{NewGun(nil, map[string]interface{}{
			"limit": 0,
		}, nil, nil, nil, nil)}
	case "Hunter":
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"speed":    1.25,
			"damage":   0.75,
			"rdir":     math.Pi / 120,
			"radius":   0.7,
			"gunbound": 0.03,
			"reload":   0.4,
			"sy":       2.23,
		}, nil, nil, nil, nil), NewGun(obj, map[string]interface{}{
			"speed":    1.25,
			"damage":   0.75,
			"rdir":     math.Pi / 120,
			"gunbound": 0.03,
			"waittime": 0.25,
			"reload":   0.4,
			"sy":       2.23,
		}, nil, nil, nil, nil)}
	case "Stalker":
		obj.Sight = 1.25
		obj.Stats = [8]int{0, 0, 0, 7, 7, 7, 5, 7}
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"speed":    1.5,
			"gunbound": 3,
			"rdir":     math.Pi / 120,
			"reload":   0.5,
			"sy":       2.45,
		}, nil, nil, nil, nil)}
		obj.Tick = func(o *Object) {
			TankTick(o)
			Invisible(o, 1.5)
		}
	case "Predator":
		obj.Sight = 1.176
		obj.Stats = [8]int{0, 0, 0, 7, 7, 7, 5, 7}
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"speed":    1.25,
			"damage":   0.75,
			"radius":   0.7,
			"rdir":     math.Pi / 120,
			"gunbound": 0.03,
			"reload":   0.4,
			"sy":       2.23,
		}, nil, nil, nil, nil), NewGun(obj, map[string]interface{}{
			"speed":    1.25,
			"damage":   0.75,
			"radius":   0.95,
			"rdir":     math.Pi / 120,
			"gunbound": 0.03,
			"waittime": 0.25,
			"reload":   0.4,
			"sy":       2.23,
		}, nil, nil, nil, nil), NewGun(obj, map[string]interface{}{
			"speed":    1.25,
			"damage":   0.75,
			"radius":   1.25,
			"rdir":     math.Pi / 120,
			"gunbound": 0.03,
			"waittime": 0.5,
			"reload":   0.4,
			"sy":       2.23,
		}, nil, nil, nil, nil)}
		obj.Tick = func(o *Object) {
			TankTick(o)
			if o.Controller != nil {
				if o.Controller.Mr && !o.IsBack {
					o.Controller.Camera.Pos = Pos{
						X: o.X + math.Cos(o.Dir)*400,
						Y: o.Y + math.Sin(o.Dir)*400,
					}
					o.IsBack = true
				}
				if !o.Controller.Mr && o.IsBack {
					o.IsBack = false
				}
			}
		}
	case "Trapper":
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"type":     "Trap",
			"health":   2,
			"reload":   0.667,
			"radius":   0.8,
			"lifetime": 24,
			"sy":       1.2,
		}, nil, nil, nil, nil)}
	case "TriTrapper":
		obj.Sight = 1.11
		var dir float64 = -math.Pi / 3 * 2
		obj.Guns = []Gun{}
		for i := 0; i < 3; i++ {
			obj.Guns = append(obj.Guns, NewGun(obj, map[string]interface{}{
				"type":     "Trap",
				"health":   2,
				"reload":   0.667,
				"radius":   0.8,
				"dir":      dir,
				"lifetime": 10,
				"sy":       1.2,
			}, nil, nil, nil, nil))
			dir += math.Pi / 3 * 2
		}
	case "Battleship":
		obj.Sight = 1.11
	case "Spike":
		obj.Sight = 1.11
		obj.Stats = [8]int{3, 10, 10, 0, 0, 0, 0, 10}
		obj.MaxStats = [8]int{10, 10, 10, 0, 0, 0, 0, 10}
		obj.Tick = func(o *Object) {
			TankTick(o)
			o.Damage += 8
		}
	case "Skimmer":
		obj.Sight = 1.11
		obj.Stats = [8]int{0, 0, 0, 5, 7, 7, 7, 7}
		var gun Gun = NewGun(obj, map[string]interface{}{
			"type":     "SkimmerBullet",
			"speed":    0.6,
			"health":   3,
			"radius":   1.75,
			"gunbound": 3,
			"bound":    0.2,
			"reload":   0.25,
			"lifetime": 4,
			"sy":       1.88,
			"isborder": false,
		}, func(o *Object) {
			if o.IsOwnCol {
				o.IsBack = o.Controller != nil && o.Controller.Mr
				o.IsOwnCol = false
			}
			if o.IsBack {
				o.Dir -= 0.05
			} else {
				o.Dir += 0.05
			}
			var dir = math.Atan2(o.Dy, o.Dx)
			o.Dx += math.Cos(dir) * o.Speed
			o.Dy += math.Sin(dir) * o.Speed

			if !o.IsDead {
				o.DeadTime -= 1000. / 60.
				if o.DeadTime <= 0 {
					o.DeadTime = -1
					o.H = 0
				}
			}
		}, nil, nil, nil)
		gun.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"health":   0.6,
			"damage":   0.6,
			"radius":   1,
			"reload":   3,
			"lifetime": 0.5,
			"autoshot": true,
			"sy":       1.8,
		}, nil, nil, nil, nil), NewGun(obj, map[string]interface{}{
			"health":   0.6,
			"damage":   0.6,
			"radius":   1,
			"reload":   3,
			"lifetime": 0.5,
			"dir":      math.Pi,
			"autoshot": true,
			"sy":       1.8,
		}, nil, nil, nil, nil)}
		obj.Guns = []Gun{
			NewGun(nil, map[string]interface{}{}, nil, nil, nil, nil),
			gun,
		}
	case "Rocketeer":
		obj.Sight = 1.11
		obj.Stats = [8]int{0, 0, 0, 7, 7, 7, 7, 5}
		var gun Gun = NewGun(obj, map[string]interface{}{
			"type":     "RocketeerBullet",
			"speed":    0.75,
			"health":   5,
			"radius":   1.333,
			"gunbound": 3,
			"bound":    0.2,
			"reload":   0.25,
			"sy":       1.88,
			"isborder": false,
		}, DefaultBulletTick, nil, nil, nil)
		gun.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"speed":    1.25,
			"health":   0.6,
			"damage":   0.3,
			"radius":   1,
			"reload":   5,
			"waittime": 8,
			"lifetime": 0.5,
			"dir":      math.Pi,
			"rdir":     math.Pi / 8,
			"autoshot": true,
			"sy":       1.8,
		}, nil, nil, nil, nil)}
		obj.Guns = []Gun{
			NewGun(nil, map[string]interface{}{}, nil, nil, nil, nil),
			gun,
		}
	default:
	}

	return obj
}

func NewAuto5() *Object {
	var obj *Object = NewObject(map[string]interface{}{
		"type":       "Auto5",
		"x":          lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
		"y":          lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
		"level":      45,
		"exp":        23536,
		"stats":      [8]int{0, 0, 0, 7, 7, 7, 5, 7},
		"maxStats":   [8]int{7, 7, 7, 7, 7, 7, 7, 7},
		"isShowName": true,
	}, TankTick, DefaultCollision, DefaultKillEvent, nil)
	var dir float64 = -math.Pi * 4 / 5
	for i := 0; i < 5; i++ {

		dir += math.Pi * 2 / 5
	}
	return obj
}

func Invisible(o *Object, t float64) {
	if o.IsShot || o.Controller == nil || o.Controller.IsMove {
		o.Opacity = math.Min(o.Opacity+0.1, 1)
	} else {
		o.Opacity = math.Max(o.Opacity-1./60./t, 0)
	}
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
