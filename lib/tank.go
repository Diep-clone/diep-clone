package lib

import (
	"math"
)

func (obj *Object) Tank() {
	obj.Type = "Tank"
	obj.Event = map[string]interface{}{
		"Tick": func(obj *Object) { // obj
			p := *obj.Owner
			if u, ok := p.(Player); ok {
				obj.Speed = (0.07 + (0.007 * obj.Variable["Stats"].([]float64)[7])) * math.Pow(0.985, obj.Variable["Level"].(float64)-1)
				obj.Damage = (20 + obj.Variable["Stats"].([]float64)[2]*4)
				obj.R = (13 * math.Pow(1.01, (obj.Variable["Level"].(float64)-1)))
				obj.Mh = (48 + obj.Variable["Level"].(float64)*2 + obj.Variable["Stats"].([]float64)[1]*20)
				obj.H += obj.Mh / 60 / 30 * (0.03 + 0.12*obj.Variable["Stats"].([]float64)[0])

				u.Camera.Z = 16 / 9 * math.Pow(0.995, obj.Variable["Level"].(float64)-1) / obj.Variable["Sight"].(float64)
			} else {
				obj.H -= obj.Mh / 60 / 10
			}
		},
		"KeyDown": nil, // obj, keyType
		"KeyPress": map[string]func(obj *Object, value interface{}){
			"moveRotate": func(obj *Object, value interface{}) {
				v := value.(float64)
				obj.Dx += math.Cos(v) * obj.Speed
				obj.Dy += math.Sin(v) * obj.Speed
			},
		}, // obj, keyType
		"KeyUp":     nil, // obj, keyType
		"Collision": nil, // obj, enemyObj
		"GetDamage": nil, // obj, enemyObj, damage
		"KillEvent": nil, // obj, enemyObj
		"DeadEvent": nil, // obj, enemyObj
	}
	obj.Variable = map[string]interface{}{
		"Level":         1,
		"MaxStats":      []int{7, 7, 7, 7, 7, 7, 7, 7},
		"Stats":         []int{0, 0, 0, 0, 0, 0, 0, 0},
		"Stat":          0,
		"Sight":         1.,
		"InvisibleTime": 0.,
	}
	obj.IsBorder = true
	obj.IsOwnCol = false
}

// Basic is god
func (obj *Object) Basic() {
	obj.Type = "Basic"
	obj.Guns = []Gun{
		*NewGun(
			obj,
			*NewObject(
				obj,
				"BasicBullet",
				"ffa",
				"",
				0,
				0,
				1,
				1,
				1,
				1,
				1,
				1,
				map[string]interface{}{
					"Tick": func(obj *Object) {
						obj.Dx += math.Cos(obj.Dir) * obj.Speed
						obj.Dy += math.Sin(obj.Dir) * obj.Speed

						obj.Dx *= 0.97
						obj.Dy *= 0.97
					},
				},
				nil,
				false,
				false,
			),
			0,
			1.88,
			0,
			math.Pi/36,
			0.1,
			0.6,
			0,
			0,
			0,
			map[string]interface{}{
				"Tick": func(gun *Gun) {
					// obj.H
				},
				"Shoot": func(gun *Gun) {

				},
			},
			nil,
			false,
			0,
		),
	}
}
