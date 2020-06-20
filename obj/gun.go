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
	IsBorder  bool    `json:"isborder"`
	IsAi      bool    `json:"isai"`
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
	IsShot    bool
	Limit     int  `json:"limit"`
	AutoShot  bool `json:"autoshot"`
}

func (obj *Object) Shot(objIndex int) {
	if obj.Guns == nil {
		return
	}
	obj.IsShot = false
	for i := 0; i < len(obj.Guns); i++ {
		g := &obj.Guns[i]
		if g.Owner == nil {
			continue
		}
		if obj.Controller != nil && obj.Controller.Ml {
			obj.IsShot = true
		}
		if g.AutoShot || obj.Controller != nil && obj.Controller.Ml {
			g.ShotTime += 1000. / 60.
			var GunOwner = obj
			for GunOwner.Owner != nil {
				GunOwner = GunOwner.Owner
			}
			if g.DelayTime <= 0 && g.WaitTime < g.ShotTime/((0.6-0.04*float64(GunOwner.Stats[6]))/g.Reload*1000) && g.Limit != 0 {
				dir := obj.Dir + g.Dir + lib.RandomRange(-g.Rdir, g.Rdir)
				var bullet Object = *NewObject(map[string]interface{}{
					"type":         g.Type,
					"team":         GunOwner.Team,
					"x":            obj.X + math.Cos(obj.Dir+g.Dir-math.Pi/2)*g.Sx*obj.R + math.Cos(obj.Dir+g.Dir)*g.Sy*obj.R,
					"y":            obj.Y + math.Sin(obj.Dir+g.Dir-math.Pi/2)*g.Sx*obj.R + math.Sin(obj.Dir+g.Dir)*g.Sy*obj.R,
					"dir":          dir,
					"speed":        (0.056 + 0.02*float64(GunOwner.Stats[3])) * g.Speed,
					"dx":           math.Cos(dir) * 4 * g.Speed,
					"dy":           math.Sin(dir) * 4 * g.Speed,
					"mh":           (8 + 6*float64(GunOwner.Stats[4])) * g.Health,
					"h":            (8 + 6*float64(GunOwner.Stats[4])) * g.Health,
					"damage":       (7 + 3*float64(GunOwner.Stats[5])) * g.Damage,
					"r":            0.4 * g.Radius * 12.9 * math.Pow(1.01, float64(obj.Level)-1),
					"bound":        g.Bound,
					"stance":       g.Stance,
					"deadtime":     1000 * g.LifeTime,
					"isShowHealth": false,
					"isTargeted":   false,
					"isBorder":     g.IsBorder,
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
				if !g.IsAi {
					bullet.SetController(GunOwner.Controller)
				}
				if g.Limit != -1 {
					g.Limit--
				}
				bullet.Guns = make([]Gun, len(g.Guns))
				copy(bullet.Guns, g.Guns) // because it is Slice!
				bullet.Owner = GunOwner
				obj.Dx -= math.Cos(obj.Dir+g.Dir) * 0.1 * g.GunBound
				obj.Dy -= math.Sin(obj.Dir+g.Dir) * 0.1 * g.GunBound
				g.DelayTime = (0.6 - 0.04*float64(GunOwner.Stats[6])) / g.Reload * 1000
				g.IsShot = true
				Objects[objIndex] = &bullet
				objIndex = len(Objects)
				Objects = append(Objects, obj)
			}
		} else {
			g.ShotTime = 0
		}
		g.DelayTime = math.Max(g.DelayTime-1000./60., 0)
	}
	return
}

// New Gun
func NewGun(own *Object, value map[string]interface{}, t func(*Object), c func(*Object, *Object), k func(*Object, *Object), d func(*Object, *Object)) Gun {
	m := map[string]interface{}{
		"type":      "Bullet",
		"speed":     1,
		"damage":    1,
		"health":    1,
		"radius":    1,
		"bound":     1,
		"stance":    0.1,
		"lifetime":  3,
		"isborder":  true,
		"isai":      false,
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
		case "Trap":
			t = DefaultTrapTick
		}
	}
	s.Tick = t
	if s.Type == "Bullet" || s.Type == "Trap" {
		s.IsBorder = false
	}
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
	return s
}
