package lib

import "math"

func TankTick(obj *Object) {
	obj.R = Grid * math.Pow(1.01, float64(obj.Level)-1)
	obj.Damage = 20 + float64(obj.Stats[2])*4
	obj.Speed = (0.07 + 0.007*float64(obj.Stats[7])) * math.Pow(0.985, float64(obj.Level)-1)

	if obj.Mh != 48.+float64(obj.Level)*2.+float64(obj.Stats[1])*20 {
		im := obj.Mh
		obj.Mh = 48. + float64(obj.Level)*2. + float64(obj.Stats[1])*20
		obj.H *= obj.Mh / im
	}

	if obj.Controller != nil {
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
