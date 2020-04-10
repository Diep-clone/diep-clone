package lib

import (
	"encoding/json"
	"math"
)

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

// Score is
type Score struct {
	Name  string
	Type  string
	Score int
}

// Scoreboard is
type Scoreboard [11]Score

// Scoreboard Push
func (sc *Scoreboard) Push(value Score) {
	var index = 0
	for ; value.Score < sc[index+1].Score && index < 10; index++ {
	}
	for j := 9; j >= index; j-- {
		sc[j+1] = sc[j]
	}
	sc[index] = value
}

// Object is
type Object struct {
	Controller   *Player `json:"controller"`
	Owner        *Object `json:"owner"`
	ID           int     `json:"id"`
	Type         string  `json:"type"`
	Team         string  `json:"team"`
	Name         string  `json:"name"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	R            float64 `json:"r"`
	Dx           float64 `json:"dx"`
	Dy           float64 `json:"dy"`
	Dir          float64 `json:"dir"`
	Level        int     `json:"level"`
	Exp          int     `json:"exp"`
	H            float64 `json:"h"`
	Mh           float64 `json:"mh"`
	Lh           float64 `json:"lh"`
	Damage       float64 `json:"damage"`
	GetDH        float64 `json:"getdh"`
	Speed        float64 `json:"speed"`
	Bound        float64 `json:"bound"`
	Stance       float64 `json:"stance"`
	Opacity      float64 `json:"opacity"`
	Sight        float64 `json:"sight"`
	Stats        [8]int  `json:"stats"`
	MaxStats     [8]int  `json:"maxStats"`
	Guns         []Gun   `json:"guns"`
	SpawnTime    int64   `json:"spawnTime"`
	HitTime      int64   `json:"hitTime"`
	DeadTime     float64 `json:"deadTime"`
	IsBorder     bool    `json:"isBorder"`
	IsOwnCol     bool    `json:"isOwnCol"`
	IsDead       bool    `json:"isDead"`
	IsCollision  bool    `json:"isCollision"`
	IsShot       bool    `json:"isShot"`
	IsShowHealth bool    `json:"isShowHealth"`
	HitObject    *Object
	Tick         func(*Object)
	Collision    func(*Object, *Object)
	KillEvent    func(*Object, *Object)
	DeadEvent    func(*Object, *Object)
}

//
func (obj *Object) ObjectTick() {
	obj.X += obj.Dx * setting.GameSpeed
	obj.Y += obj.Dy * setting.GameSpeed

	obj.Dx *= 0.97 //math.Pow(0.97, setting.GameSpeed)
	obj.Dy *= 0.97 //math.Pow(0.97, setting.GameSpeed)

	if obj.H <= 0 {
		obj.IsDead = true
		obj.H = 0
		if obj.HitObject != nil && obj.HitObject.KillEvent != nil {
			obj.HitObject.KillEvent(obj.HitObject, obj)
		}
	}

	if obj.Controller != nil && now()-int64(30000.*setting.GameSpeed) > obj.HitTime {
		obj.H += obj.Mh / 60 / 10 * setting.GameSpeed
	}

	obj.H += obj.GetDH * setting.GameSpeed
	obj.GetDH = 0

	if obj.H > obj.Mh {
		obj.H = obj.Mh
	}

	obj.Lh = obj.H

	obj.IsCollision = false
}

//
func DefaultCollision(a *Object, b *Object) {
	dir := math.Atan2(a.Y-b.Y, a.X-b.X)

	a.Dx += math.Cos(dir) * math.Min(b.Bound*a.Stance, 6)
	a.Dy += math.Sin(dir) * math.Min(b.Bound*a.Stance, 6)

	if a.Team == b.Team {
		return
	}

	a.HitTime = now()

	if b.Lh-a.Damage*setting.GameSpeed <= 0 {
		a.H -= b.Damage * (b.Lh / (a.Damage * setting.GameSpeed)) * setting.GameSpeed
	} else {
		a.H -= b.Damage * setting.GameSpeed
	}

	a.HitObject = b
	a.IsCollision = true
}

var objID int = 1

//
func (o Object) SocketObj() map[string]interface{} {
	var ownerID int = -1
	if o.Owner != nil {
		ownerID = o.Owner.ID
	}
	if !o.IsShowHealth {
		o.Mh = 1
		o.H = 1
	}
	return map[string]interface{}{
		"id":      o.ID,
		"team":    o.Team,
		"type":    o.Type,
		"x":       floor(o.X, 2),
		"y":       floor(o.Y, 2),
		"r":       floor(o.R, 1),
		"dir":     floor(o.Dir, 2),
		"mh":      floor(o.Mh, 1),
		"h":       floor(o.H, 1),
		"opacity": floor(o.Opacity, 2),
		"exp":     o.Exp,
		"name":    o.Name,
		"owner":   ownerID,
		"isDead":  o.IsDead,
	}
}

func (o *Object) SetController(p *Player) {
	o.Controller = p
}

func NewObject(value map[string]interface{}, guns []Gun, t func(*Object), c func(*Object, *Object), k func(*Object, *Object), d func(*Object, *Object)) *Object {
	m := map[string]interface{}{
		"id":           objID,
		"type":         "Square",
		"team":         "ffa",
		"x":            -999999,
		"y":            -999999,
		"r":            Grid,
		"level":        1,
		"h":            10,
		"mh":           10,
		"lh":           10,
		"damage":       8,
		"speed":        0.07,
		"bound":        1,
		"stance":       1,
		"opacity":      1,
		"sight":        1,
		"spawnTime":    now(),
		"hitTime":      now(),
		"deadTime":     -1,
		"isBorder":     true,
		"isOwnCol":     true,
		"isDead":       false,
		"isCollision":  false,
		"isShowHealth": true,
	}
	objID++

	for key, val := range value {
		m[key] = val
	}

	jsonString, _ := json.Marshal(m)

	s := Object{}
	json.Unmarshal(jsonString, &s)
	s.Guns = guns
	s.Tick = t
	s.Collision = c
	s.KillEvent = k
	s.DeadEvent = d
	return &s
}
