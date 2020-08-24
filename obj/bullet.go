package obj

import (
	"app/lib"
	"math"
)

// TODO : Should fix object AI

func DefaultTrapTick(obj *Object) {
	if lib.Now()-obj.SpawnTime > 3000 {
		obj.IsOwnCol = false
	} else {
		obj.IsOwnCol = true
	}

	if obj.IsShowName == false {
		var GunOwner = obj.Owner
		for GunOwner.Owner != nil {
			GunOwner = GunOwner.Owner
		}
		obj.Dx += obj.Dx * 0.1 * float64(GunOwner.Stats[3])
		obj.Dy += obj.Dy * 0.1 * float64(GunOwner.Stats[3])
		obj.IsShowName = true
	}

	if !obj.IsDead {
		obj.DeadTime -= 1000. / 60.
		if obj.DeadTime <= 0 {
			obj.DeadTime = -1
			obj.H = 0
		}
	}
}

func DefaultBulletTick(obj *Object) {
	obj.IsOwnCol = false
	obj.Dx += math.Cos(obj.Dir) * obj.Speed
	obj.Dy += math.Sin(obj.Dir) * obj.Speed

	if !obj.IsDead {
		obj.DeadTime -= 1000. / 60.
		if obj.DeadTime <= 0 {
			obj.DeadTime = -1
			obj.H = 0
		}
	}
}

//
func DefaultHealCollision(a *Object, b *Object) {
	if a.Team == b.Team && b.Controller != nil {
		b.H += 5 + (1.5 * float64(a.Stats[5]))
		a.H -= (5 + (1.5 * float64(a.Stats[5]))) * 30

		if b.Mh < b.H {
			a.H += (b.H - b.Mh) * 30
			b.H = b.Mh
		}

		if a.H <= 0 {
			a.IsCollision = true
		}
	} else {
		DefaultCollision(a, b)
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
	var targetdis float64
	if target != nil {
		targetdis = lib.Distance(target.X, target.Y, obj.Owner.X, obj.Owner.Y)
	}

	if obj.IsBack {
		if dis < 60 {
			obj.IsBack = false
		} else if dis < 550 && target != nil && targetdis < 550 {
			obj.IsBack = false
			obj.Target = target
		} else {
			obj.Dir = math.Atan2(obj.Owner.Y-obj.Y, obj.Owner.X-obj.X)
			obj.Dx += math.Cos(obj.Dir) * obj.Speed
			obj.Dy += math.Sin(obj.Dir) * obj.Speed
		}
	} else if obj.Target != nil {
		if (obj.Target.IsDead || obj.Owner == obj.Target || obj.Target.Team == obj.Team || !obj.Target.IsTargeted) && targetdis < 550 {
			obj.Target = target
		} else if dis > 550 || targetdis >= 550 {
			obj.Target = nil
		} else {
			obj.Dir = math.Atan2(obj.Target.Y-obj.Y, obj.Target.X-obj.X)
			obj.Dx += math.Cos(obj.Dir) * obj.Speed
			obj.Dy += math.Sin(obj.Dir) * obj.Speed
		}
	} else if dis < lib.RandomRange(130, 150) {
		var dir float64 = math.Atan2(obj.Owner.Y-obj.Y, obj.Owner.X-obj.X) - math.Pi/2
		obj.Dx += math.Cos(dir) * 0.02
		obj.Dy += math.Sin(dir) * 0.02
		obj.Dir = math.Atan2(obj.Dy, obj.Dx)
		if targetdis < 550 {
			obj.Target = target
		}
	} else {
		obj.IsBack = true
	}
}

func BattleshipDroneAi(obj *Object) {
	obj.IsOwnCol = false
	if obj.Owner.DeadTime == 0 {
		obj.H = 0
		return
	}

	if !obj.IsDead {
		obj.DeadTime -= 1000. / 60.
		if obj.DeadTime <= 0 {
			obj.DeadTime = -1
			obj.H = 0
		}
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
	var target *Object = NearObject(obj, 1000, 0, math.Pi)

	if obj.Target != nil {
		if obj.Target.IsDead || obj.Owner == obj.Target || obj.Target.Team == obj.Team || !obj.Target.IsTargeted {
			obj.Target = target
		} else {
			obj.Dir = math.Atan2(obj.Target.Y-obj.Y, obj.Target.X-obj.X)
			obj.Dx += math.Cos(obj.Dir) * obj.Speed
			obj.Dy += math.Sin(obj.Dir) * obj.Speed
		}
	} else {
		obj.Dx += math.Cos(obj.Dir) * obj.Speed
		obj.Dy += math.Sin(obj.Dir) * obj.Speed
		obj.Target = target
	}
}
