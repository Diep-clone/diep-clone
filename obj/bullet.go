package obj

import "math"

func DefaultBulletTick(obj *Object) {
	obj.IsOwnCol = false
	obj.Dx += math.Cos(obj.Dir) * obj.Speed
	obj.Dy += math.Sin(obj.Dir) * obj.Speed

	obj.DeadTime -= 1000. / 60.
	if obj.DeadTime <= 0 {
		obj.DeadTime = -1
		obj.H = 0
	}
}
