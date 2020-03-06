package lib

import "time"

// Camera is
type Camera struct {
	x float64
	y float64
	z float64
}

// Player is
type Player struct {
	Id            string
	IsMove        bool
	Mx            float64
	My            float64
	Camera        Camera
	Keys          map[string]float64
	StartTime     int64
	ControlObject *Object
}

// NewPlayer is make player
func NewPlayer(id string) *Player {
	p := Player{}
	p.Id = id
	p.IsMove = false
	p.Mx = 0
	p.My = 0
	p.Camera = Camera{
		x: 0,
		y: 0,
		z: 1,
	}
	p.Keys = map[string]float64{}
	p.StartTime = time.Now().Unix()

	return &p
}

// Player's mouse point set
func (p *Player) SetMousePoint(x float64, y float64) {
	p.Mx = x
	p.My = y
}

// Player's key set
func (p *Player) SetKey(key string, value float64) {
	p.Keys[key] = value
	//obj := p.ControlObject
	//ok := obj != nil
	switch key {
	case "moveRotate":
		p.IsMove = value > 0
		/*case "1":
			if ok { obj.Variable["HealthRegen"]}
		case "2":
			if ok { obj.Variable["MaxHealth"]}
		case "3":
			if ok { obj.Variable["BodyDamage"]}
		case "4":
			if ok { obj.Variable["BulletSpeed"]}
		case "5":
			if ok { obj.Variable["BulletHealth"]}
		case "6":
			if ok { obj.Variable["BulletDamage"]}
		case "7":
			if ok { obj.Variable["Reload"]}
		case "8":
			if ok { obj.Variable["Movement"]}*/

	}
}

// Object is
type Object struct {
	Owner     *interface{}
	Id        int
	Paths     []interface{}
	Color     [3]int
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
	Guns      []Gun
	Event     map[string]interface{}
	Variable  map[string]interface{}
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
	paths []interface{},
	color [3]int,
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
	o.Paths = paths
	o.Color = color
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
	o.Guns = []Gun{}
	o.Event = event
	o.Variable = variable
	o.HitTime = time.Now().Unix()
	o.DeadTime = -1
	o.HitObject = &o
	o.IsBorder = isBorder
	o.IsOwnCol = isOwnCol
	o.IsDead = false

	return &o
}

// Gun is
type Gun struct {
	Object    Object
	Paths     []interface{}
	Color     [3]int
	Sx        float64
	Sy        float64
	Dir       float64
	Rdir      float64
	Reload    float64 // default delay
	WaitTime  float64 // first wait time
	DelayTime float64 // when shot
	ShotTime  float64 // click time to delay
	AutoShot  bool
	Limit     int
	Bullets   []Object
}
