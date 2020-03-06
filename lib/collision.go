package lib

import (
	"math"
	"time"
)

// CollisionEvent is
func CollisionEvent(a *Object, b *Object) bool {
	dir := math.Atan2(a.y-b.y, a.x-b.x)

	if a == *b.owner || b == *a.owner {
		return false
	}

	a.dx += math.Cos(dir) * math.Min(b.bound*a.stance, 6)
	a.dy += math.Sin(dir) * math.Min(b.bound*a.stance, 6)
	b.dx -= math.Cos(dir) * math.Min(a.bound*b.stance, 6)
	b.dy -= math.Sin(dir) * math.Min(a.bound*b.stance, 6)

	if a.team != "-1" && b.team != "-1" && a.team == b.team {
		return false
	}

	a.hitTime = time.Now().Unix()
	b.hitTime = time.Now().Unix()

	a.hitObject = b
	b.hitObject = a

	if b.lh-a.damage <= 0 {
		a.h -= b.damage * (b.lh / a.damage)
	} else {
		a.h -= b.damage
	}
	if a.lh-b.damage <= 0 {
		b.h -= a.damage * (a.lh / b.damage)
	} else {
		b.h -= a.damage
	}
	if a.h < 0 {
		a.h = 0
	}
	if b.h < 0 {
		b.h = 0
	}

	return true
}

// MaxObj is Object Max length
var MaxObj = 5

// Area is
type Area struct {
	x  float64
	y  float64
	x2 float64
	y2 float64
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
		nodes:   []Quadtree{},
	}
	return &q
}

// Split is
func (q Quadtree) Split() {
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

// Getindex is
func (q Quadtree) Getindex(area interface{}) int {
	obj, b := area.(Object)
	x := q.x + q.w/2
	y := q.y + q.h/2
	if b {
		if obj.x+obj.r >= x && obj.x-obj.r <= x || obj.y+obj.r >= y && obj.y-obj.r <= y {
			return -1
		}
	} else {
		obj, _ := area.(Area)
		if obj.x >= x && obj.x2 <= x || obj.y >= y && obj.y2 <= y {
			return -1
		}
	}
	// 2 1
	// 3 4
	if obj.x > x {
		if obj.y > y {
			return 4
		}
		return 1
	}

	if obj.y > y {
		return 3
	}

	return 2
}

// Insert insert quadtree
func (q Quadtree) Insert(obj *Object) {
	var i, index = 0, -1

	// obj가 this의 관할이 아닐 때
	if obj.x+obj.r < q.x || obj.x-obj.r > q.x+q.w ||
		obj.y+obj.r < q.y || obj.y-obj.r > q.y+q.h {
		return
	}

	// 자식 노드가 있을 때
	if q.nodes != nil {
		index = q.Getindex(obj)
		if index != -1 {
			q.nodes[index].Insert(obj)
			return
		}
	}

	q.objects = append(q.objects, obj)

	if len(q.objects) > MaxObj {
		if q.nodes != nil {
			q.Split()
		}

		for i < len(q.objects) {
			index = q.Getindex(q.objects[i])
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
	index := q.Getindex(area)
	var returnObject []*Object

	if o, ok := area.(Object); ok {
		for _, obj := range q.objects {
			if !obj.isDead && (obj.owner != o.owner || obj.isOwnCol && o.isOwnCol) {
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
