package obj

import (
	"math"

	"app/lib"
)

func TankTick(obj *Object) {
	if obj.IsDead {
		return
	}
	obj.R = 12.9 * math.Pow(1.01, float64(obj.Level)-1)
	obj.Damage = 20 + float64(obj.Stats[2])*4
	obj.Speed = (0.07 + 0.007*float64(obj.Stats[7])) * math.Pow(0.985, float64(obj.Level)-1)

	if obj.Mh != 48.+float64(obj.Level)*2.+float64(obj.Stats[1])*20 {
		im := obj.Mh
		obj.Mh = 48. + float64(obj.Level)*2. + float64(obj.Stats[1])*20
		obj.H *= obj.Mh / im
	}

	if obj.Controller != nil && obj.H > 0 {
		obj.H += obj.Mh / 60 / 30 * (0.03 + 0.12*float64(obj.Stats[0]))

		for i := 0; i < 8; i++ {
			if obj.Stats[i] > obj.MaxStats[i] {
				obj.Controller.Stat += obj.Stats[i] - obj.MaxStats[i]
				obj.Stats[i] = obj.MaxStats[i]
			}
		}
	} else {
		obj.H -= obj.Mh / 60 / 10
	}
}

// TODO : ChangeTank Method
func (o *Object) ChangeTank(c *Object) {
	o.Type = c.Type
	o.Bounce = c.Bounce
	o.Stance = c.Stance
	o.Sight = c.Sight
	o.MaxStats = c.MaxStats
	o.Guns = c.Guns
	o.Tick = c.Tick
	o.Collision = c.Collision
	o.KillEvent = c.KillEvent
	o.DeadEvent = c.DeadEvent
}

// TODO : other Tank
func NewTank(t string) *Object {
	var obj *Object = NewObject(map[string]interface{}{
		"type":       t,
		"x":          lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
		"y":          lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
		"level":      45,
		"exp":        23536,
		"stats":      [8]int{0, 0, 0, 5, 7, 7, 7, 7},
		"maxStats":   [8]int{7, 7, 7, 7, 7, 7, 7, 7},
		"stance":     0.3,
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
		obj.Sight = 1.11
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
		var dir float64 = -math.Pi / 3. * 2.
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
			dir += math.Pi / 3. * 2.
		}
	case "Auto5":
		obj.Sight = 1.11
		for i := 0.; i < 5.; i++ {
			obj.SubObjects = append(obj.SubObjects, AutoGun(obj, map[string]interface{}{
				"speed":  1.2,
				"damage": 0.4,
				"radius": 1.2,
				"bound":  0.333,
			}, math.Pi*2./5.*i, math.Pi/2.))
		}
	case "Auto3":
		obj.Sight = 1.11
		for i := 0.; i < 3.; i++ {
			obj.SubObjects = append(obj.SubObjects, AutoGun(obj, map[string]interface{}{
				"speed":  1.2,
				"damage": 0.4,
				"radius": 1.2,
				"bound":  0.333,
			}, math.Pi*2./3.*i, math.Pi/2.))
		}
	case "AutoTrapper":
		obj.Sight = 1.11
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"type":     "Trap",
			"health":   2,
			"reload":   0.667,
			"radius":   0.8,
			"lifetime": 24,
			"sy":       1.2,
		}, nil, nil, nil, nil)}
		obj.SubObjects = append(obj.SubObjects, nil)
		obj.SubObjects = append(obj.SubObjects, AutoGun(obj, map[string]interface{}{
			"speed":  1.2,
			"damage": 0.3,
			"radius": 1.2,
			"bound":  0.125,
		}, 0, 0))
	case "Battleship":
		obj.Sight = 1.11
		obj.Stats = [8]int{0, 0, 0, 7, 7, 7, 7, 5}
		var r []float64 = []float64{math.Pi / 2, -math.Pi / 2}
		var posx []float64 = []float64{0.25, -0.25}
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				obj.Guns = append(obj.Guns, NewGun(obj, map[string]interface{}{
					"type":     "Drone",
					"damage":   0.15,
					"radius":   0.5,
					"dir":      r[i],
					"bound":    3.4,
					"lifetime": 4,
					"isai":     (j == 1),
					"sx":       posx[j],
					"sy":       1.5,
					"isborder": false,
				}, BattleshipDroneAi, nil, nil, nil))
			}
		}
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
		// and Other tanks was Making by users
	case "Dispersion":
		for i := 0; i < 7; i++ {
			obj.Guns = append(obj.Guns, NewGun(obj, map[string]interface{}{
				"health": 0.7,
				"damage": 0.5,
				"rdir":   math.Pi / 9,
			}, nil, nil, nil, nil))
		}
	case "Diffusion":
		obj.Stats = [8]int{0, 0, 0, 0, 7, 7, 7, 5}
		var gunlist []Gun
		for i := 0; i < 8; i++ {
			gunlist = append(gunlist, NewGun(obj, map[string]interface{}{
				"damage":   0.4,
				"reload":   0,
				"lifetime": 2,
				"sy":       0,
				"rdir":     math.Pi,
				"autoshot": true,
			}, nil, nil, nil, nil))
		}
		var gun Gun = NewGun(obj, map[string]interface{}{
			"health": 2,
			"damage": 1.75,
			"reload": 0.25,
			"radius": 1.75,
			"bound":  3,
		}, nil, nil, nil, func(a *Object, b *Object) {
			a.Guns = gunlist
			a.Shot()
		})
		obj.Guns = []Gun{gun}
	case "Deception":
		obj.Sight = 1.11
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"type":     "Trap",
			"health":   2,
			"reload":   0.667,
			"radius":   0.8,
			"lifetime": 24,
			"sy":       1.2,
		}, func(o *Object) {
			if lib.Now()-o.SpawnTime > 3000 {
				o.IsOwnCol = false
				o.Controller = nil
				Invisible(o, 5)
			} else {
				o.IsOwnCol = true
			}

			if o.IsShowName == false {
				var GunOwner = o.Owner
				for GunOwner.Owner != nil {
					GunOwner = GunOwner.Owner
				}
				o.Dx += o.Dx * 0.1 * float64(GunOwner.Stats[3])
				o.Dy += o.Dy * 0.1 * float64(GunOwner.Stats[3])
				o.IsShowName = true
			}

			if !o.IsDead {
				o.DeadTime -= 1000. / 60.
				if o.DeadTime <= 0 {
					o.DeadTime = -1
					o.H = 0
				}
			}
		}, func(a *Object, b *Object) {
			DefaultCollision(a, b)
			a.Opacity = math.Min(a.Opacity+0.1, 1)
		}, nil, nil)}
	case "Dropper":
		obj.Sight = 1.11
		obj.Stats = [8]int{0, 0, 0, 7, 7, 7, 7, 5}
		for i := 0.; i < 2; i++ {
			var gun Gun = NewGun(obj, map[string]interface{}{
				"type":     "DropperBullet",
				"speed":    0.55,
				"damage":   0.7,
				"health":   2,
				"reload":   0.167,
				"autoshot": true,
				"gunbound": 3,
				"bound":    0.5,
				"sy":       1.6,
				"lifetime": 0,
				"dir":      math.Pi/2 - i*math.Pi,
				"limit":    3,
			}, DefaultDroneTick, nil, nil, nil)
			gun.Guns = []Gun{NewGun(obj, map[string]interface{}{
				"type":     "Trap",
				"speed":    0.01,
				"health":   2,
				"radius":   0.8,
				"reload":   0.333,
				"waittime": 8,
				"lifetime": 8,
				"dir":      math.Pi,
				"autoshot": true,
				"sy":       1.2,
			}, nil, nil, nil, nil)}
			obj.Guns = append(obj.Guns, gun)
		}
		obj.SubObjects = append(obj.SubObjects, nil)
		var oim = NewObject(map[string]interface{}{
			"team":     "nil",
			"type":     "Gun",
			"mh":       0,
			"damage":   0,
			"isBorder": false,
			"isOwnCol": false,
		}, func(o *Object) {
			o.R = o.Owner.R * 0.5
			o.X = o.Owner.X
			o.Y = o.Owner.Y
			if o.Owner.H == 0 || o.Owner.IsDead == true {
				o.IsDead = true
			}
		}, nil, nil, nil)
		oim.Owner = obj
		obj.SubObjects = append(obj.SubObjects, oim)
	case "Shielder":
		obj.Stats = [8]int{0, 7, 0, 0, 7, 7, 7, 5}
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{}, nil, nil, nil, nil)}
		var oim = NewObject(map[string]interface{}{
			"type":     "Shield",
			"team":     nil,
			"level":    1,
			"opacity":  0,
			"stance":   0,
			"isOwnCol": false,
			"isDead":   true,
			"deadTime": 0,
		}, DefaultShieldTick, DefaultCollision, DefaultKillEvent, nil)
		oim.Owner = obj
		obj.SubObjects = []*Object{nil, oim}
	case "Mechanic":
		obj.Stats = [8]int{0, 7, 0, 5, 7, 7, 7, 0}
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"damage":   0.5,
			"radius":   0.5,
			"gunbound": 0.4,
			"sx":       0.325,
			"sy":       2.1,
		}, nil, nil, nil, nil), NewGun(obj, map[string]interface{}{
			"damage":   0.5,
			"radius":   0.5,
			"gunbound": 0.4,
			"waittime": 0.5,
			"sx":       -0.325,
			"sy":       2.1,
		}, nil, nil, nil, nil)}
		var oim = NewObject(map[string]interface{}{
			"type":     "MechanicArmor",
			"team":     nil,
			"level":    1,
			"opacity":  0,
			"stance":   0,
			"isOwnCol": false,
			"isDead":   true,
			"deadTime": 0,
		}, DefaultShieldTick, DefaultCollision, DefaultKillEvent, nil)
		oim.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"speed":    0.75,
			"damage":   2.5,
			"health":   2,
			"radius":   2,
			"gunbound": 15,
			"bound":    0.1,
			"reload":   0.25,
			"sy":       1.48,
		}, nil, nil, nil, nil)}
		oim.Owner = obj
		obj.SubObjects = []*Object{nil, oim}
	case "Cure":
		obj.Stats = [8]int{0, 0, 0, 7, 7, 7, 7, 5}
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"bound":  0,
			"reload": 0.85,
		}, nil, DefaultHealCollision, nil, nil)}
		var oim = NewObject(map[string]interface{}{
			"team":     "nil",
			"type":     "Cross",
			"mh":       0,
			"r":        1,
			"damage":   0,
			"isBorder": false,
			"isOwnCol": false,
		}, func(o *Object) {
			o.R = o.Owner.R * 0.5
			o.X = o.Owner.X
			o.Y = o.Owner.Y
			o.Dir += 0.005
			if o.Owner.H == 0 || o.Owner.IsDead == true {
				o.IsDead = true
			}
		}, nil, nil, nil)
		oim.Owner = obj
		obj.SubObjects = []*Object{nil, oim}
	case "Follower":
		obj.Sight = 1.11
		obj.Stats = [8]int{0, 0, 0, 7, 7, 7, 7, 5}
		obj.Guns = []Gun{
			NewGun(obj, map[string]interface{}{
				"type":     "Square",
				"speed":    0.6,
				"damage":   0.7,
				"health":   2,
				"radius":   1,
				"gunbound": 0.4,
				"bound":    0.2,
				"reload":   0.667,
				"lifetime": 0,
				"sy":       1.2,
				"isborder": true,
				"limit":    8,
			}, func(o *Object) {
				if o.Owner.DeadTime == 0 {
					o.H = 0
					return
				}
				o.Dir += 0.05
				if o.IsBack {
					if o.Controller != nil && o.Controller.Mr {
						o.Sight = math.Atan2((o.Owner.Y+o.Controller.My)-o.Y, (o.Owner.X+o.Controller.Mx)-o.X)
					}
					o.Dx += math.Cos(o.Sight) * o.Speed
					o.Dy += math.Sin(o.Sight) * o.Speed
				} else if o.Controller != nil && o.Controller.Mr {
					o.Sight = math.Atan2((o.Owner.Y+o.Controller.My)-o.Y, (o.Owner.X+o.Controller.Mx)-o.X)
					o.IsBack = true
				}
			}, nil, nil, nil),
		}
	case "Lifesteal":
		obj.Stats = [8]int{0, 0, 0, 7, 7, 7, 5, 7}
		obj.Sight = 1.11
		obj.Guns = []Gun{NewGun(obj, map[string]interface{}{
			"type":     "Drone",
			"speed":    0.6,
			"damage":   0.7,
			"health":   2,
			"gunbound": 0.4,
			"reload":   0.333,
			"autoshot": true,
			"sy":       1.6,
			"limit":    8,
		}, nil, nil, nil, func(a *Object, b *Object) {
			a.Owner.H += 10
		})}
		obj.SubObjects = append(obj.SubObjects, nil)
		var oim = NewObject(map[string]interface{}{
			"team":     "nil",
			"type":     "TriangleGun",
			"mh":       0,
			"r":        1,
			"damage":   0,
			"isBorder": false,
			"isOwnCol": false,
		}, func(o *Object) {
			o.R = o.Owner.R * 0.5
			o.X = o.Owner.X
			o.Y = o.Owner.Y
			o.Dir = o.Owner.Dir

			if o.Owner.H == 0 || o.Owner.IsDead == true {
				o.IsDead = true
			}
		}, nil, nil, nil)
		oim.Owner = obj
		obj.SubObjects = append(obj.SubObjects, oim)
	default:
	}

	return obj
}

func AutoGun(owner *Object, gunValue map[string]interface{}, dir, rdir float64) *Object {
	var o *Object = NewObject(map[string]interface{}{
		"team":     "nil",
		"type":     "AutoGun",
		"speed":    rdir,
		"sight":    dir,
		"mh":       0,
		"r":        1,
		"damage":   0,
		"isBorder": false,
		"isOwnCol": false,
	}, func(obj *Object) {
		obj.R = obj.Owner.R * 0.5
		obj.Owner.Dx += obj.Dx
		obj.Owner.Dy += obj.Dy
		obj.Dx = 0
		obj.Dy = 0

		if obj.Owner.H == 0 || obj.Owner.IsDead == true {
			obj.IsDead = true
		}

		if obj.Speed == 0. {
			obj.X = obj.Owner.X
			obj.Y = obj.Owner.Y
			obj.Dir += 0.005

			var target *Object = NearObject(obj, 400, 0, math.Pi)

			if obj.Target == nil {
				obj.Target = target
			} else {
				if (obj.Target.IsDead || obj.Owner == obj.Target || obj.Target.Team == obj.Team || !obj.Target.IsTargeted) && lib.Distance(obj.Target.X, obj.Target.Y, obj.X, obj.Y) > 500 {
					obj.Target = target
				} else if obj.Target != nil && lib.Distance(obj.Target.X, obj.Target.Y, obj.X, obj.Y) > 550 {
					obj.Target = nil
				}
			}

			if obj.Target != nil {
				obj.Guns[0].AutoShot = true
				obj.Dir = math.Atan2(obj.Target.Y-obj.Y, obj.Target.X-obj.X)
			} else {
				obj.Guns[0].AutoShot = false
			}
		} else {
			obj.X = obj.Owner.X + math.Cos(obj.Sight)*obj.Owner.R*0.8
			obj.Y = obj.Owner.Y + math.Sin(obj.Sight)*obj.Owner.R*0.8
			if obj.Target == nil {
				obj.Dir = obj.Sight
			}
			obj.Sight = obj.Sight + 0.005
			if obj.Sight > math.Pi*2 {
				obj.Sight -= math.Pi * 2
			}

			var target *Object = NearObject(obj, 400, obj.Sight, obj.Speed)

			if obj.Owner.Controller != nil && obj.Owner.Controller.Mr && lib.DirDis(obj.Sight, math.Atan2(obj.Y-(obj.Owner.Y+obj.Owner.Controller.My), obj.X-(obj.Owner.X+obj.Owner.Controller.Mx))) < rdir {
				obj.Guns[0].AutoShot = true
				obj.Dir = math.Atan2(obj.Y-(obj.Owner.Y+obj.Owner.Controller.My), obj.X-(obj.Owner.X+obj.Owner.Controller.Mx))
			} else if obj.Owner.Controller != nil && obj.Owner.Controller.Ml && lib.DirDis(obj.Sight, math.Atan2((obj.Owner.Y+obj.Owner.Controller.My)-obj.Y, (obj.Owner.X+obj.Owner.Controller.Mx)-obj.X)) < rdir {
				obj.Guns[0].AutoShot = true
				obj.Dir = math.Atan2((obj.Owner.Y+obj.Owner.Controller.My)-obj.Y, (obj.Owner.X+obj.Owner.Controller.Mx)-obj.X)
			} else {
				if obj.Target == nil {
					obj.Target = target
				} else {
					if (obj.Target.IsDead || obj.Owner == obj.Target || obj.Target.Team == obj.Team || !obj.Target.IsTargeted) && lib.Distance(obj.Target.X, obj.Target.Y, obj.X, obj.Y) > 500 {
						obj.Target = target
					} else if obj.Target != nil && lib.Distance(obj.Target.X, obj.Target.Y, obj.X, obj.Y) > 550 {
						obj.Target = nil
					}
				}

				if obj.Target != nil {
					obj.Guns[0].AutoShot = true
					obj.Dir = math.Atan2(obj.Target.Y-obj.Y, obj.Target.X-obj.X)
				} else {
					obj.Guns[0].AutoShot = false
				}
			}
		}
	}, nil, DefaultKillEvent, nil)
	o.Owner = owner
	o.Guns = []Gun{NewGun(o, gunValue, nil, nil, nil, nil)}
	return o
}

func DefaultShieldTick(o *Object) {
	o.X = o.Owner.X
	o.Y = o.Owner.Y
	o.Owner.Dx += o.Dx
	o.Owner.Dy += o.Dy
	o.Dx = 0
	o.Dy = 0
	o.Damage = o.Owner.Damage
	o.Mh = o.Owner.Mh * (0.3 + 0.1*float64(o.Owner.Stats[4]))
	o.Team = o.Owner.Team
	o.Dir = o.Owner.Dir
	o.SetController(o.Owner.Controller)

	if o.Owner.H == 0 || o.Owner.IsDead == true {
		o.DeadTime = 99999
		o.IsDead = true
	}

	if o.IsDead {
		if o.DeadTime == -1 {
			o.DeadTime = 25000 - float64(o.Owner.Stats[6])*3000
		} else if o.DeadTime <= 0 {
			o.IsDead = false
			o.DeadTime = -1
			o.H = o.Mh
		} else {
			o.Opacity = math.Max(o.Opacity-0.05, 0)
			if o.Opacity > 0 {
				o.R += o.R * 0.02
			}
			o.DeadTime = math.Max(o.DeadTime-1000./60., 0.)
		}

		return
	}

	o.R = o.Owner.R + (float64(o.Level) * 5.)
	o.Opacity = 0.5

	for _, e := range Qt.Retrieve(Area{
		X: o.X - o.R,
		Y: o.Y - o.R,
		W: o.R * 2,
		H: o.R * 2,
	}) {
		if o != e && !o.IsDead && !e.IsDead && !e.IsCollision && (e.Owner != o.Owner || e.IsOwnCol && o.IsOwnCol) && o != e.Owner && e != o.Owner && !lib.Contains(e.HitObjects, o) {
			if (o.X-e.X)*(o.X-e.X)+(o.Y-e.Y)*(o.Y-e.Y) < (o.R+e.R)*(o.R+e.R) {
				if o.Collision != nil {
					o.Collision(o, e)
				}
				if e.Collision != nil {
					e.Collision(e, o)
				}
			}
		}
	}
	o.Owner.HitObjects = append(o.Owner.HitObjects, o.HitObjects...)
	o.HitObjects = nil
}

func Invisible(o *Object, t float64) {
	if !o.IsDead && (o.IsShot || o.Controller != nil) && o.Controller.IsMove {
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
