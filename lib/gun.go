package lib

// Gun is
type Gun struct {
	Owner     *Object
	Object    Object
	Sx        float64
	Sy        float64
	Dir       float64
	Rdir      float64
	Bound     float64
	Reload    float64 // default delay
	WaitTime  float64 // first wait time
	DelayTime float64 // when shot
	ShotTime  float64 // click time to delay
	Event     map[string]interface{}
	Variable  map[string]interface{}
	AutoShot  bool
	Limit     int
	Bullets   []Object
}

// New Gun1!!!!!!!!!!11!!!
func NewGun(
	owner *Object,
	obj Object,
	sx float64,
	sy float64,
	dir float64,
	rdir float64,
	bound float64,
	reload float64,
	waitTime float64,
	delayTime float64,
	shotTime float64,
	event map[string]interface{},
	variable map[string]interface{},
	autoShot bool,
	limit int) *Gun {
	g := Gun{}
	g.Owner = owner
	g.Object = obj
	g.Sx = sx
	g.Sy = sy
	g.Dir = dir
	g.Rdir = rdir
	g.Bound = bound
	g.Reload = reload
	g.WaitTime = waitTime
	g.DelayTime = delayTime
	g.ShotTime = shotTime
	g.Event = event
	g.Variable = variable
	g.AutoShot = autoShot
	g.Limit = limit
	g.Bullets = []Object{}
	return &g
}
