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
	id            string
	moveRotate    float64
	isMove        bool
	mx            float64
	my            float64
	camera        Camera
	keys          map[string]float64
	startTime     int64
	controlObject *Object
}

// NewPlayer is make player
func NewPlayer(id string) *Player {
	p := Player{}
	p.id = id
	p.moveRotate = 0
	p.isMove = false
	p.mx = 0
	p.my = 0
	p.camera = Camera{
		x: 0,
		y: 0,
		z: 1,
	}

	p.keys = map[string]float64{}
	p.startTime = time.Now().Unix()

	return &p
}

// Object is
type Object struct {
	owner     *interface{}
	id        int
	shapeType int
	color     [3]int
	team      string
	name      string
	x         float64
	y         float64
	dx        float64
	dy        float64
	r         float64
	dir       float64
	exp       int
	h         float64
	mh        float64
	lh        float64
	damage    float64
	speed     float64
	bound     float64
	stance    float64
	opacity   float64
	guns      []Gun
	event     map[string]interface{}
	variable  map[string]interface{}
	hitTime   int64
	deadTime  int64
	hitObject *Object
	isBorder  bool
	isOwnCol  bool
	isDead    bool
}

var objID int = 0

// NewObject is
func NewObject(
	own interface{},
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
	o.owner = &own
	o.id = objID
	objID++
	o.team = team
	o.name = name
	o.x = x
	o.y = y
	o.dx = 0
	o.dy = 0
	o.mh = h
	o.h = h
	o.lh = h
	o.damage = da
	o.speed = sp
	o.bound = bo
	o.stance = st
	o.guns = []Gun{}
	o.event = event
	o.variable = variable
	o.hitTime = time.Now().Unix()
	o.deadTime = -1
	o.hitObject = &o
	o.isBorder = isBorder
	o.isOwnCol = isOwnCol
	o.isDead = false

	return &o
}

// Gun is
type Gun struct {
	object    Object
	color     [3]int
	paths     []interface{}
	sx        float64
	sy        float64
	dir       float64
	rdir      float64
	reload    float64 // default delay
	waitTime  float64 // first wait time
	delayTime float64 // when shot
	shotTime  float64 // click time to delay
	autoShot  bool
	bullets   []Object
}
