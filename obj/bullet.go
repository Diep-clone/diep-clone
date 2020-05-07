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

func DefaultDroneTick(obj *Object) {
	p := obj.Controller
	if p == nil {
		return
	}
	if obj.Owner.DeadTime == 0 {
		obj.H = 0
		return
	}
	if p.Mr {
		obj.Dir = math.Atan2(obj.Y-(obj.Owner.Y+p.My), obj.X-(obj.Owner.X+p.Mx))
		obj.Dx += math.Cos(obj.Dir) * obj.Speed
		obj.Dy += math.Sin(obj.Dir) * obj.Speed
	} else if p.Ml {
		obj.Dir = math.Atan2((obj.Owner.Y+p.My)-obj.Y, (obj.Owner.X+p.Mx)-obj.X)
		obj.Dx += math.Cos(obj.Dir) * obj.Speed
		obj.Dy += math.Sin(obj.Dir) * obj.Speed
	}
}
