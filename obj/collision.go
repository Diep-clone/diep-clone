package obj

import (
	"app/lib"
	"math"
)

// MaxObj is Object Max length
var MaxObj = 4

var Qt Quadtree

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
func NewQuadtree(x, y, w, h float64) *Quadtree {
	var q Quadtree = Quadtree{
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		level:   1,
		objects: []*Object{},
		nodes:   nil,
	}
	return &q
}

// split is
func (q *Quadtree) split() {
	xx := [4]float64{0, 1, 0, 1}
	yy := [4]float64{0, 0, 1, 1}
	q.nodes = make([]Quadtree, 4)
	for i := 0; i < 4; i++ {
		q.nodes[i] = Quadtree{
			x:       q.x + xx[i]*q.w/2,
			y:       q.y + yy[i]*q.h/2,
			w:       q.w / 2,
			h:       q.h / 2,
			level:   q.level + 1,
			objects: []*Object{},
			nodes:   nil,
		}
	}
}

// getIndex is
func (q Quadtree) getIndex(area Area) int {
	x := q.x + q.w/2
	y := q.y + q.h/2

	if area.X <= x && area.X+area.W >= x || area.Y <= y && area.Y+area.H >= y {
		return -1
	}
	// 1 0
	// 2 3
	if area.X > x {
		if area.Y > y {
			return 3
		}
		return 0
	}

	if area.Y > y {
		return 2
	}
	return 1
}

// Insert insert quadtree
func (q *Quadtree) Insert(obj *Object) {
	if q.level == 1 && (obj.X+obj.R < q.x || obj.X-obj.R > q.x+q.w || obj.Y+obj.R < q.y || obj.Y-obj.R > q.y+q.h) {
		return
	}

	var index = -1

	if q.nodes != nil {
		index = q.getIndex(Area{
			X: obj.X - obj.R,
			Y: obj.Y - obj.R,
			W: obj.R * 2,
			H: obj.R * 2,
		})
		if index != -1 {
			q.nodes[index].Insert(obj)
			return
		}
	}

	q.objects = append(q.objects, obj)

	if len(q.objects) > MaxObj {
		if q.nodes == nil {
			q.split()
		}

		for i := 0; i < len(q.objects); {
			index = q.getIndex(Area{
				X: q.objects[i].X - q.objects[i].R,
				Y: q.objects[i].Y - q.objects[i].R,
				W: q.objects[i].R * 2,
				H: q.objects[i].R * 2,
			})
			if index != -1 {
				q.nodes[index].Insert(q.objects[i])
				q.objects[i] = q.objects[len(q.objects)-1]
				q.objects = q.objects[:len(q.objects)-1]
			} else {
				i++
			}
		}
	}
}

// RetrieveArea is
func (q Quadtree) Retrieve(area Area) []*Object {
	var index int = q.getIndex(area)
	var returnObjects []*Object = q.objects

	if q.nodes != nil {
		if index != -1 {
			for _, obj := range q.nodes[index].Retrieve(area) {
				returnObjects = append(returnObjects, obj)
			}
		} else {
			for i := 0; i < 4; i++ {
				for _, obj := range q.nodes[i].Retrieve(area) {
					returnObjects = append(returnObjects, obj)
				}
			}
		}
	}

	return returnObjects
}

func NearObject(o *Object, ran, dir, rdir float64) *Object {
	var objList []*Object = Qt.Retrieve(Area{
		X: o.X - ran,
		Y: o.Y - ran,
		W: ran * 2,
		H: ran * 2,
	})
	var nearObj *Object = nil

	for _, obj := range objList {
		if math.Abs(dir-math.Atan2(o.Y-obj.Y, o.X-obj.X)) < rdir && lib.Distance(o.X, o.Y, obj.X, obj.Y) < ran && o != obj && o.Owner != obj && o.Team != obj.Team && !obj.IsDead {
			if nearObj == nil || lib.Distance(o.X, o.Y, obj.X, obj.Y) < lib.Distance(o.X, o.Y, nearObj.X, nearObj.Y) {
				nearObj = obj
			}
		}
	}

	return nearObj
}

// Clear is
func (q *Quadtree) Clear() {
	q.objects = nil

	if q.nodes != nil {
		for i := 0; i < 4; i++ {
			q.nodes[i].Clear()
		}
	}

	q.nodes = nil
}
