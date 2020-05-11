package obj

import (
	"app/lib"
	"math"
)

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
	if obj.Owner.DeadTime == 0 {
		obj.H = 0
		return
	}
	p := obj.Controller
	if p != nil {
		if p.Mr {
			obj.Dir = math.Atan2(obj.Y-(obj.Owner.Y+p.My), obj.X-(obj.Owner.X+p.Mx))
			obj.Dx += math.Cos(obj.Dir) * obj.Speed
			obj.Dy += math.Sin(obj.Dir) * obj.Speed
			obj.IsBack = false
			return
		} else if p.Ml {
			obj.Dir = math.Atan2((obj.Owner.Y+p.My)-obj.Y, (obj.Owner.X+p.Mx)-obj.X)
			obj.Dx += math.Cos(obj.Dir) * obj.Speed
			obj.Dy += math.Sin(obj.Dir) * obj.Speed
			obj.IsBack = false
			return
		}
	}
	var dis float64 = lib.Distance(obj.Owner.X, obj.Owner.Y, obj.X, obj.Y)
	var target *Object = NearObject(obj, 300, 0, math.Pi)

	if obj.IsBack {
		if dis < 300 {
			obj.IsBack = false
		} else {
			obj.Dir = math.Atan2(obj.Owner.Y-obj.Y, obj.Owner.X-obj.X)
			obj.Dx += math.Cos(obj.Dir) * obj.Speed
			obj.Dy += math.Sin(obj.Dir) * obj.Speed
		}
	} else if obj.Target != nil {
		if obj.Target.IsDead || obj.Owner == obj.Target || obj.Target.Team == obj.Team {
			obj.Target = target
		} else if dis > 700 {
			obj.Target = nil
		} else {
			obj.Dir = math.Atan2(obj.Target.Y-obj.Y, obj.Target.X-obj.X)
			obj.Dx += math.Cos(obj.Dir) * obj.Speed
			obj.Dy += math.Sin(obj.Dir) * obj.Speed
		}
	} else if dis < lib.RandomRange(139.32, 150) {
		var dir float64 = math.Atan2(obj.Owner.Y-obj.Y, obj.Owner.X-obj.X) - math.Pi/2
		obj.Dx += math.Cos(dir) * 0.02
		obj.Dy += math.Sin(dir) * 0.02
		obj.Dir = math.Atan2(obj.Dy, obj.Dx)
		if target != nil {
			obj.Target = target
		}
	} else {
		obj.IsBack = true
	}
}
