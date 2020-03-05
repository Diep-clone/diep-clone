package lib

import "time"

// Camera is
type Camera struct {
	x float32
	y float32
	z float32
}

// Player is
type Player struct {
	id            string
	moveRotate    float32
	isMove        bool
	mx            float32
	my            float32
	camera        Camera
	keys          map[string]float32
	startTime     int64
	controlObject *Object
}

func newPlayer(id string) *Player {
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

	p.keys = map[string]float32{}
	p.startTime = time.Now().Unix()

	return &p
}

// Object is
type Object struct {
	owner     interface{}
	id        int
	shapeType int
	color     [3]int
	team      string
	name      string
	x         float32
	y         float32
	dx        float32
	dy        float32
	r         float32
	dir       float32
	exp       int
	h         float32
	mh        float32
	damage    float32
	speed     float32
	bound     float32
	stance    float32
	opacity   float32
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

func newObject(
	own interface{},
	team string,
	name string,
	x float32,
	y float32,
	r float32,
	h float32,
	da float32,
	sp float32,
	bo float32,
	st float32,
	event map[string]interface{},
	variable map[string]interface{},
	isBorder bool,
	isOwnCol bool) *Object {
	o := Object{}
	o.owner = own
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
	sx        float32
	sy        float32
	dir       float32
	rdir      float32
	reload    float32 // default delay
	waitTime  float32 // first wait time
	delayTime float32 // when shot
	shotTime  float32 // click time to delay
	autoShot  bool
	bullets   []Object
}
