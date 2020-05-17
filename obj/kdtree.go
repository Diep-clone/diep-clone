package obj

type K2Dtree struct {
	root *node
}

//New KDtree is
func NewK2Dtree() K2Dtree {
	return K2Dtree{}
}

//KDtree insert
func (kd *K2Dtree) Insert(obj *Object) {
	if kd.root == nil {
		kd.root = &node{
			x:       obj.X,
			y:       obj.Y,
			objects: []*Object{obj},
		}
	} else {
		kd.root.Insert(obj, 1)
	}
}

type node struct {
	x       float64
	y       float64
	left    *node
	right   *node
	objects []*Object
}

func (n *node) Insert(obj *Object, level int) {
	var b bool
	if level&1 == 1 {
		b = obj.X < n.x
	} else {
		b = obj.Y < n.y
	}
	if b {
		if n.left == nil {
			n.left = &node{}
		} else {
			n.left.Insert(obj, (level+1)%2)
		}
	} else {
		if n.right == nil {
			n.right = &node{}
		} else {
			n.right.Insert(obj, (level+1)%2)
		}
	}
}
