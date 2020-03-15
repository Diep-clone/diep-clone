package lib

import (
	"math"
	"time"
)

// CollisionEvent is
func CollisionEvent(a *Object, b *Object) bool {
	dir := math.Atan2(a.C.Pos.Y-b.C.Pos.Y, a.C.Pos.Y-b.C.Pos.Y)

	if a == *b.Owner || b == *a.Owner {
		return false
	}

	if f, ok := a.Event["Collision"].(func(a *Object, b *Object)); ok {
		f(a, b)
	}
	if f, ok := b.Event["Collision"].(func(a *Object, b *Object)); ok {
		f(b, a)
	}

	a.Dx += math.Cos(dir) * math.Min(b.Bound*a.Stance, 6)
	a.Dy += math.Sin(dir) * math.Min(b.Bound*a.Stance, 6)
	b.Dx -= math.Cos(dir) * math.Min(a.Bound*b.Stance, 6)
	b.Dy -= math.Sin(dir) * math.Min(a.Bound*b.Stance, 6)

	if a.Team != "ffa" && b.Team != "ffa" && a.Team == b.Team {
		return false
	}

	a.HitTime = time.Now().Unix()
	b.HitTime = time.Now().Unix()

	a.HitObject = b
	b.HitObject = a

	if b.Lh-a.Damage <= 0 {
		a.H -= b.Damage * (b.Lh / a.Damage)
	} else {
		a.H -= b.Damage
	}
	if a.Lh-b.Damage <= 0 {
		b.H -= a.Damage * (a.Lh / b.Damage)
	} else {
		b.H -= a.Damage
	}
	if f, ok := a.Event["GetDamage"].(func(a *Object, b *Object)); ok {
		f(a, b)
	}
	if f, ok := b.Event["GetDamage"].(func(a *Object, b *Object)); ok {
		f(b, a)
	}
	if a.H < 0 {
		a.H = 0
		if f, ok := a.Event["DeadEvent"].(func(a *Object, b *Object)); ok {
			f(a, b)
		}
		if f, ok := b.Event["KillEvent"].(func(a *Object, b *Object)); ok {
			f(b, a)
		}
	}
	if b.H < 0 {
		b.H = 0
		if f, ok := a.Event["DeadEvent"].(func(a *Object, b *Object)); ok {
			f(b, a)
		}
		if f, ok := b.Event["KillEvent"].(func(a *Object, b *Object)); ok {
			f(a, b)
		}
	}

	return true
}

// MaxObj is Object Max length
var MaxObj = 5

// Area is
type Area struct {
	X  float64
	Y  float64
	X2 float64
	Y2 float64
}

// Quadtree is
type Quadtree struct {
	x       float64
	y       float64
	w       float64
	h       float64
	level   int
	objects []*Object
	nodes   []Quadtree
}

// NewQuadtree is
func NewQuadtree(x float64, y float64, w float64, h float64, level int) *Quadtree {
	q := Quadtree{
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		level:   level,
		objects: []*Object{},
		nodes:   nil,
	}
	return &q
}

// split is
func (q Quadtree) split() {
	xx := [4]float64{0, 1, 0, 1}
	yy := [4]float64{0, 0, 1, 1}
	for i := 0; i < 4; i++ {
		q.nodes[i] = Quadtree{
			x:       q.x + xx[i]*q.w/2,
			y:       q.y + yy[i]*q.h/2,
			w:       q.w / 2,
			h:       q.h / 2,
			level:   q.level + 1,
			objects: []*Object{},
			nodes:   []Quadtree{},
		}
	}
}

// getIndex is
func (q Quadtree) getIndex(area interface{}) int {
	obj, b := area.(Object)
	x := q.x + q.w/2
	y := q.y + q.h/2
	if b {
		if obj.C.Pos.X+obj.C.R >= x && obj.C.Pos.X-obj.C.R <= x || obj.C.Pos.Y+obj.C.R >= y && obj.C.Pos.Y-obj.C.R <= y {
			return -1
		}
	} else {
		obj, _ := area.(Area)
		if obj.X >= x && obj.X2 <= x || obj.Y >= y && obj.Y2 <= y {
			return -1
		}
	}
	// 2 1
	// 3 4
	if obj.C.Pos.X > x {
		if obj.C.Pos.Y > y {
			return 4
		}
		return 1
	}

	if obj.C.Pos.Y > y {
		return 3
	}

	return 2
}

// Insert insert quadtree
func (q Quadtree) Insert(obj *Object) {
	var i, index = 0, -1

	// obj가 this의 관할이 아닐 때
	if obj.C.Pos.X+obj.C.R < q.x || obj.C.Pos.X-obj.C.R > q.x+q.w ||
		obj.C.Pos.Y+obj.C.R < q.y || obj.C.Pos.Y-obj.C.R > q.y+q.h {
		return
	}

	// 자식 노드가 있을 때
	if q.nodes != nil {
		index = q.getIndex(obj)
		if index != -1 {
			q.nodes[index].Insert(obj)
			return
		}
	}

	q.objects = append(q.objects, obj)

	if len(q.objects) > MaxObj {
		if q.nodes != nil {
			q.split()
		}

		for i < len(q.objects) {
			index = q.getIndex(q.objects[i])
			if index != -1 {
				q.nodes[index].Insert(q.objects[i])
				q.objects = q.objects[i:]
			} else {
				i++
			}
		}
	}
}

// Retrieve is
func (q Quadtree) Retrieve(area interface{}) []*Object {
	index := q.getIndex(area)
	var returnObject []*Object

	if o, ok := area.(Object); ok {
		for _, obj := range q.objects {
			if !obj.IsDead && (obj.Owner != o.Owner || obj.IsOwnCol && o.IsOwnCol) {
				returnObject = append(returnObject, obj)
			}
		}

	} else {
		returnObject = q.objects
	}

	if q.nodes != nil {
		if index != -1 {
			for _, obj := range q.nodes[index].Retrieve(area) {
				returnObject = append(returnObject, obj)
			}
		} else {
			for i := 0; i < 4; i++ {
				for _, obj := range q.nodes[i].Retrieve(area) {
					returnObject = append(returnObject, obj)
				}
			}
		}
	}

	return returnObject
}

// Clear is
func (q Quadtree) Clear() {
	q.objects = nil

	if q.nodes != nil {
		for i := 0; i < 4; i++ {
			q.nodes[i].Clear()
		}
	}

	q.nodes = nil
}
