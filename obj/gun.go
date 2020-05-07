package obj

import (
	"encoding/json"
	"math"

	"app/lib"
)

// Gun is
type Gun struct {
	Owner     *Object
	Type      string  `json:"type"`
	Speed     float64 `json:"speed"`
	Damage    float64 `json:"damage"`
	Health    float64 `json:"health"`
	Radius    float64 `json:"radius"`
	Bound     float64 `json:"bound"`
	Stance    float64 `json:"stance"`
	LifeTime  float64 `json:"lifetime"`
	Guns      []Gun
	Tick      func(*Object)
	Collision func(*Object, *Object)
	KillEvent func(*Object, *Object)
	DeadEvent func(*Object, *Object)

	Sx        float64 `json:"sx"`
	Sy        float64 `json:"sy"`
	Dir       float64 `json:"dir"`
	Rdir      float64 `json:"rdir"`
	GunBound  float64 `json:"gunbound"`
	Reload    float64 `json:"reload"`    // default delay
	DelayTime float64 `json:"delaytime"` // when shot
	WaitTime  float64 `json:"waittime"`  // first wait time
	ShotTime  float64 `json:"shottime"`  // click time to delay
	Limit     int     `json:"limit"`
	AutoShot  bool    `json:"autoshot"`
}

func (g *Gun) Shot() {
	if g.Owner == nil || g.Owner.Controller == nil {
		return
	}
	if g.AutoShot || g.Owner.Controller.Ml {
		g.ShotTime += 1000. / 60.
		var GunOwner = g.Owner
		for GunOwner.Owner != nil {
			GunOwner = GunOwner.Owner
		}
		g.Owner.IsShot = g.Owner.Controller.Ml
		if g.DelayTime <= 0 && g.WaitTime < g.ShotTime/((0.6-0.04*float64(GunOwner.Stats[6]))/g.Reload*1000) && g.Limit != 0 {
			dir := g.Owner.Dir + g.Dir + lib.RandomRange(-g.Rdir, g.Rdir)
			var bullet Object = *NewObject(map[string]interface{}{
				"type":         g.Type,
				"team":         GunOwner.Team,
				"x":            g.Owner.X + math.Cos(g.Owner.Dir+g.Dir-math.Pi/2)*g.Sx*g.Owner.R + math.Cos(g.Owner.Dir+g.Dir)*g.Sy*g.Owner.R,
				"y":            g.Owner.Y + math.Sin(g.Owner.Dir+g.Dir-math.Pi/2)*g.Sx*g.Owner.R + math.Sin(g.Owner.Dir+g.Dir)*g.Sy*g.Owner.R,
				"dir":          dir,
				"speed":        (0.056 + 0.02*float64(GunOwner.Stats[3])) * g.Speed,
				"dx":           math.Cos(dir) * (4 + 0.4*float64(GunOwner.Stats[3])) * g.Speed,
				"dy":           math.Sin(dir) * (4 + 0.4*float64(GunOwner.Stats[3])) * g.Speed,
				"mh":           (8 + 6*float64(GunOwner.Stats[4])) * g.Health,
				"h":            (8 + 6*float64(GunOwner.Stats[4])) * g.Health,
				"damage":       (7 + 3*float64(GunOwner.Stats[5])) * g.Damage,
				"r":            0.4 * g.Radius * 12.9 * math.Pow(1.01, float64(g.Owner.Level)-1),
				"bound":        g.Bound,
				"stance":       g.Stance,
				"deadtime":     1000 * g.LifeTime,
				"isShowHealth": false,
			}, func(o *Object) {
				if g.Limit == -1 {
					o.Tick = g.Tick
				}
				if o.IsDead {
					g.Limit++
					o.Tick = g.Tick
				}
				g.Tick(o)
			}, g.Collision, g.KillEvent, g.DeadEvent)
			bullet.SetController(GunOwner.Controller)
			if g.Limit != -1 {
				g.Limit--
			}
			bullet.Guns = g.Guns
			bullet.Owner = GunOwner
			g.Owner.Dx -= math.Cos(g.Owner.Dir+g.Dir) * 0.1 * g.Bound
			g.Owner.Dy -= math.Sin(g.Owner.Dir+g.Dir) * 0.1 * g.Bound
			g.DelayTime = (0.6 - 0.04*float64(GunOwner.Stats[6])) / g.Reload * 1000
			Objects = append(Objects, &bullet)
		}
	} else {
		g.ShotTime = 0
	}
	g.DelayTime = math.Max(g.DelayTime-1000./60., 0)
	return
}

// New Gun1!!!!!!!!!!11!!!
func NewGun(own *Object, value map[string]interface{}, t func(*Object), c func(*Object, *Object), k func(*Object, *Object), d func(*Object, *Object)) *Gun {
	m := map[string]interface{}{
		"type":      "Bullet",
		"speed":     1,
		"damage":    1,
		"health":    1,
		"radius":    1,
		"bound":     1,
		"stance":    1,
		"lifetime":  3,
		"sx":        0,
		"sy":        1.88,
		"dir":       0,
		"rdir":      math.Pi / 36,
		"gunbound":  1,
		"reload":    1,
		"delaytime": 0,
		"waittime":  0,
		"shottime":  0,
		"limit":     -1,
		"autoshot":  false,
	}

	for key, val := range value {
		m[key] = val
	}

	jsonString, _ := json.Marshal(m)

	var s Gun
	json.Unmarshal(jsonString, &s)
	s.Owner = own
	if t == nil {
		switch s.Type {
		case "Bullet":
			t = DefaultBulletTick
		case "Drone":
			t = DefaultDroneTick
		}
	}
	s.Tick = t
	if s.Type == "Drone" && value["lifetime"] == nil {
		s.LifeTime = -0.001
	}
	if c == nil {
		c = DefaultCollision
	}
	s.Collision = c
	if k == nil {
		k = DefaultKillEvent
	}
	s.KillEvent = k
	s.DeadEvent = d
	return &s
}
