package obj

import "math"

func DefaultBulletTick(obj *Object) {
	obj.Dx += math.Cos(obj.Dir) * obj.Speed
	obj.Dy += math.Sin(obj.Dir) * obj.Speed

	obj.Guns[0].DelayTime -= 1000. / 60.
	if obj.Guns[0].DelayTime <= 0 {
		obj.H = 0
	}
}
