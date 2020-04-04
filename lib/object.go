package lib

import (
	"encoding/json"
	"math"
	"time"
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
	Color int
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
	Controller  *Player `json:"controller"`
	Owner       *Object `json:"owner"`
	ID          int     `json:"id"`
	Type        string  `json:"type"`
	Color       int     `json:"color"`
	Team        string  `json:"team"`
	Name        string  `json:"name"`
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	R           float64 `json:"r"`
	Dx          float64 `json:"dx"`
	Dy          float64 `json:"dy"`
	Dir         float64 `json:"dir"`
	Level       int     `json:"level"`
	Exp         int     `json:"exp"`
	H           float64 `json:"h"`
	Mh          float64 `json:"mh"`
	Lh          float64 `json:"lh"`
	Damage      float64 `json:"damage"`
	Speed       float64 `json:"speed"`
	Bound       float64 `json:"bound"`
	Stance      float64 `json:"stance"`
	Opacity     float64 `json:"opacity"`
	Sight       float64 `json:"sight"`
	MaxStats    [8]int  `json:"maxStats"`
	Guns        []Gun   `json:"guns"`
	SpawnTime   int64   `json:"spawnTime"`
	HitTime     int64   `json:"hitTime"`
	DeadTime    int64   `json:"deadTime"`
	IsBorder    bool    `json:"isBorder"`
	IsOwnCol    bool    `json:"isOwnCol"`
	IsDead      bool    `json:"isDead"`
	IsCollision bool    `json:"isCollision"`
}

//
func (obj *Object) ObjectTick() {
	obj.X += obj.Dx
	obj.Y += obj.Dy

	obj.Dx *= 0.97
	obj.Dy *= 0.97

	if obj.IsDead {
		return
	}

	if obj.H <= 0 {
		obj.IsDead = true
		obj.H = 0
	}

	if time.Now().Unix()-30000 > obj.HitTime {
		obj.H += obj.Mh / 60 / 10
	}

	imMh := obj.Mh

	if obj.H > obj.Mh {
		obj.H = obj.Mh
	}

	if obj.Mh != imMh {
		obj.H *= obj.Mh / imMh
	}
}

//
func (o *Object) Collision(dir float64, b float64, team string, d float64, h float64) {
	o.Dx += math.Cos(dir) * math.Min(b*o.Stance, 6)
	o.Dy += math.Sin(dir) * math.Min(b*o.Stance, 6)

	if team != "ffa" && o.Team == team {
		return
	}

	o.HitTime = time.Now().Unix()

	if h-o.Damage <= 0 {
		o.H -= d * (h / o.Damage)
	} else {
		o.H -= d
	}

	o.IsCollision = true
}

var objID int = 0

//
func (o Object) SocketObj() map[string]interface{} {
	var ownerID int = -1
	if o.Owner != nil {
		ownerID = o.Owner.ID
	}
	return map[string]interface{}{
		"id":          o.ID,
		"type":        o.Type,
		"color":       o.Color,
		"x":           floor(o.X, 2),
		"y":           floor(o.Y, 2),
		"r":           floor(o.R, 1),
		"dir":         floor(o.Dir, 2),
		"maxhealth":   floor(o.Mh, 1),
		"health":      floor(o.H, 1),
		"opacity":     floor(o.Opacity, 2),
		"exp":         o.Exp,
		"name":        o.Name,
		"owner":       ownerID,
		"isDead":      o.IsDead,
		"isCollision": o.IsCollision,
	}
}

func NewObject(value map[string]interface{}) *Object {
	m := map[string]interface{}{
		"ID":          objID,
		"Type":        "Circle",
		"Color":       2,
		"Team":        "FFA",
		"X":           -999999,
		"Y":           -999999,
		"R":           Grid,
		"Level":       1,
		"H":           10,
		"Mh":          10,
		"Lh":          10,
		"Damage":      8,
		"Speed":       0.07,
		"Bound":       1,
		"Stance":      1,
		"Opacity":     1,
		"Sight":       1,
		"SpawnTime":   time.Now().Unix(),
		"HitTime":     time.Now().Unix(),
		"DeadTime":    -1,
		"IsBorder":    true,
		"IsOwnCol":    true,
		"IsDead":      false,
		"IsCollision": false,
	}
	objID++

	for key, val := range value {
		m[key] = val
	}

	jsonString, _ := json.Marshal(m)

	s := Object{}
	json.Unmarshal(jsonString, &s)
	return &s
}
