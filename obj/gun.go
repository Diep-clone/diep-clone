package obj

import (
	"encoding/json"
	"log"
	"math"

	"app/lib"
)

// Gun is
type Gun struct {
	Owner     *Object
	Bullet    Object
	Sx        float64 `json:"sx"`
	Sy        float64 `json:"sy"`
	Dir       float64 `json:"dir"`
	Rdir      float64 `json:"rdir"`
	Bound     float64 `json:"bound"`
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
	log.Println("AA")
	if g.AutoShot || g.Owner.Controller.Ml {
		g.ShotTime += 1000 / 60
		var GunOwner = g.Owner
		for GunOwner.Owner != nil {
			GunOwner = GunOwner.Owner
		}
		log.Println("Bb")
		g.Owner.IsShot = true
		if g.DelayTime <= 0 && g.WaitTime <= g.ShotTime/((0.6-0.04*float64(GunOwner.Stats[6]))/g.Reload*1000) {
			log.Println("Ww")
			dir := g.Owner.Dir + g.Dir + lib.RandomRange(-g.Rdir, g.Rdir)
			bullet := g.Bullet
			bullet.ID = objID
			objID++
			bullet.Type = "Bullet"
			bullet.Team = GunOwner.Team
			bullet.Owner = GunOwner
			bullet.X = g.Owner.X + math.Cos(g.Owner.Dir+g.Dir-math.Pi/2)*g.Sx*g.Owner.R + math.Cos(g.Owner.Dir+g.Dir)*g.Sy*g.Owner.R
			bullet.Y = g.Owner.Y + math.Sin(g.Owner.Dir+g.Dir-math.Pi/2)*g.Sx*g.Owner.R + math.Sin(g.Owner.Dir+g.Dir)*g.Sy*g.Owner.R
			bullet.Dir = dir
			bullet.Speed = (0.056 + 0.02*float64(GunOwner.Stats[3])) * bullet.Speed
			bullet.Dx = math.Cos(dir) * bullet.Speed * 20
			bullet.Dx = math.Sin(dir) * bullet.Speed * 20
			bullet.Mh = (8 + 6*float64(GunOwner.Stats[4])) * bullet.Mh
			bullet.H = (8 + 6*float64(GunOwner.Stats[4])) * bullet.H
			bullet.Damage = (7 + 3*float64(GunOwner.Stats[5])) * bullet.Damage
			bullet.R = 10 * bullet.R * math.Pow(1.01, float64(g.Owner.Level)-1)
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
func NewGun(value map[string]interface{}, own *Object, b Object) *Gun {
	m := map[string]interface{}{
		"sx":        0,
		"sy":        1.88,
		"dir":       0,
		"rdir":      math.Pi / 36,
		"bound":     1,
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

	s := Gun{}
	json.Unmarshal(jsonString, &s)
	s.Owner = own
	s.Bullet = b
	return &s
}
