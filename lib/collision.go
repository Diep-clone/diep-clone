package lib

// MaxObj is Object Max length
var MaxObj = 4

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
func (q *Quadtree) split() {
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
		if obj.X+obj.R >= x && obj.X-obj.R <= x || obj.Y+obj.R >= y && obj.Y-obj.R <= y {
			return -1
		}
	} else {
		obj, _ := area.(Area)
		if obj.X <= x && obj.X+obj.W >= x || obj.Y <= y && obj.Y+obj.H >= y {
			return -1
		}
	}
	// 2 1
	// 3 4
	if obj.X > x {
		if obj.Y > y {
			return 4
		}
		return 1
	}

	if obj.Y > y {
		return 3
	}

	return 2
}

// Insert insert quadtree
func (q *Quadtree) Insert(obj *Object) {
	var index = -1

	if obj.X+obj.R < q.x || obj.X-obj.R > q.x+q.w ||
		obj.Y+obj.R < q.y || obj.Y-obj.R > q.y+q.h {
		return
	}

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

		for i := 0; i < len(q.objects); {
			index = q.getIndex(q.objects[i])
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

// Retrieve is
func (q Quadtree) Retrieve(area interface{}) []*Object {
	index := q.getIndex(area)
	var returnObjects []*Object

	if o, ok := area.(*Object); ok {
		for _, obj := range q.objects {
			if !obj.IsDead && (obj.Owner != o.Owner || obj.IsOwnCol && o.IsOwnCol) && o != obj.Owner && obj != o.Owner {
				returnObjects = append(returnObjects, obj)
			}
		}
	} else {
		returnObjects = q.objects
	}

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
