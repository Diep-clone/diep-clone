package obj

import (
	"encoding/json"
	"math"

	"app/lib"
)

// Gun is
type Gun struct {
	Owner     *Object `json:"owner"`
	Bullet    Object  `json:"bullet"`
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

func (g *Gun) Shot() *Object {
	if g.AutoShot || g.Owner.Controller.Ml {
		g.ShotTime += 1000 / 60
		var GunOwner = g.Owner
		for GunOwner.Owner != nil {
			GunOwner = GunOwner.Owner
		}
		g.Owner.IsShot = true
		if g.DelayTime <= 0 && g.WaitTime <= g.ShotTime/((0.6-0.04*float64(GunOwner.Stats[6]))/g.Reload*1000) {
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
			Objects = append(Objects, &bullet)
		}
	}
	return nil
}

// New Gun1!!!!!!!!!!11!!!
func NewGun(value map[string]interface{}) *Gun {
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
	return &s
}
