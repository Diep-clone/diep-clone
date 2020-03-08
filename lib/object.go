package lib

import "time"

// Pos is
type Pos struct {
	X float64
	Y float64
}

// Circle is
type Circle struct {
	Pos Pos
	R   float64
}

// Object is
type Object struct {
	Owner     *interface{}
	Id        int
	Type      string
	Team      string
	Name      string
	X         float64
	Y         float64
	Dx        float64
	Dy        float64
	R         float64
	Dir       float64
	Exp       int
	H         float64
	Mh        float64
	Lh        float64
	Damage    float64
	Speed     float64
	Bound     float64
	Stance    float64
	Opacity   float64
	Guns      []*Gun
	Event     map[string]interface{}
	Variable  map[string]interface{}
	SpawnTime int64
	HitTime   int64
	DeadTime  int64
	HitObject *Object
	IsBorder  bool
	IsOwnCol  bool
	IsDead    bool
}

var objID int = 0

// NewObject is
func NewObject(
	own interface{},
	t string,
	team string,
	name string,
	x float64,
	y float64,
	r float64, // radius
	h float64, // health
	da float64, // damage
	sp float64, // speed
	bo float64, // bound
	st float64, // stance
	event map[string]interface{},
	variable map[string]interface{},
	isBorder bool,
	isOwnCol bool) *Object {
	o := Object{}
	o.Owner = &own
	o.Id = objID
	objID++
	o.Type = t
	o.Team = team
	o.Name = name
	o.X = x
	o.Y = y
	o.Dx = 0
	o.Dy = 0
	o.Mh = h
	o.H = h
	o.Lh = h
	o.Damage = da
	o.Speed = sp
	o.Bound = bo
	o.Stance = st
	o.Guns = []*Gun{}
	o.Event = event
	o.Variable = variable
	o.SpawnTime = time.Now().Unix()
	o.HitTime = time.Now().Unix()
	o.DeadTime = -1
	o.HitObject = &o
	o.IsBorder = isBorder
	o.IsOwnCol = isOwnCol
	o.IsDead = false

	return &o
}
