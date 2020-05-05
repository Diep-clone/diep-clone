package obj

import (
	"app/lib"
	"encoding/binary"
	"encoding/json"
	"math"
)

var Objects []*Object

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
	IsShowName   bool    `json:"isShowName"`
	HitObject    *Object
	Tick         func(*Object)
	Collision    func(*Object, *Object)
	KillEvent    func(*Object, *Object)
	DeadEvent    func(*Object, *Object)
}

//
func (obj *Object) ObjectTick() {
	obj.X += obj.Dx * lib.GameSetting.GameSpeed
	obj.Y += obj.Dy * lib.GameSetting.GameSpeed

	obj.Dx *= 1 - 0.03*lib.GameSetting.GameSpeed //math.Pow(0.97, setting.GameSpeed)
	obj.Dy *= 1 - 0.03*lib.GameSetting.GameSpeed //math.Pow(0.97, setting.GameSpeed)

	if obj.IsBorder { // 화면 밖으로 벗어나는가?
		if obj.X > lib.GameSetting.MapSize.X+lib.Grid*4 {
			obj.X = lib.GameSetting.MapSize.X + lib.Grid*4
		}
		if obj.X < -lib.GameSetting.MapSize.X-lib.Grid*4 {
			obj.X = -lib.GameSetting.MapSize.X - lib.Grid*4
		}
		if obj.Y > lib.GameSetting.MapSize.Y+lib.Grid*4 {
			obj.Y = lib.GameSetting.MapSize.Y + lib.Grid*4
		}
		if obj.Y < -lib.GameSetting.MapSize.Y-lib.Grid*4 {
			obj.Y = -lib.GameSetting.MapSize.Y - lib.Grid*4
		}
	}

	if obj.H <= 0 && !obj.IsDead {
		obj.IsDead = true
		obj.H = 0
		return
	}

	if obj.Guns != nil {
		for i := 0; i < len(obj.Guns); i++ {
			obj.Guns[i].Shot()
		}
	}

	if obj.Controller != nil && lib.Now()-int64(30000.*lib.GameSetting.GameSpeed) > obj.HitTime {
		obj.H += obj.Mh / 60 / 10 * lib.GameSetting.GameSpeed
	}

	obj.H += obj.GetDH * lib.GameSetting.GameSpeed
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

	a.HitTime = lib.Now()

	if b.Lh-a.Damage*lib.GameSetting.GameSpeed <= 0 {
		a.H -= b.Damage * (b.Lh / (a.Damage * lib.GameSetting.GameSpeed)) * lib.GameSetting.GameSpeed
	} else {
		a.H -= b.Damage * lib.GameSetting.GameSpeed
	}

	if a.H <= 0 {
		if b.KillEvent != nil {
			b.KillEvent(b, a)
		}
	}

	a.HitObject = b
	a.IsCollision = true
}

func DefaultKillEvent(o *Object, deader *Object) {
	var owner *Object = o
	for owner.Owner != nil {
		owner = owner.Owner
	}
	owner.Exp += int(math.Min(float64(deader.Exp), 23536.))
}

var objID int = 1

//
func (o Object) ObjBinaryData() []byte {
	var data []byte = make([]byte, 64)
	if !o.IsShowHealth {
		o.Mh = 1
		o.H = 1
	}
	binary.BigEndian.PutUint32(data[0:4], uint32(o.ID))
	binary.BigEndian.PutUint64(data[4:12], math.Float64bits(lib.Floor(o.X, 2)))
	binary.BigEndian.PutUint64(data[12:20], math.Float64bits(lib.Floor(o.Y, 2)))
	binary.BigEndian.PutUint64(data[20:28], math.Float64bits(lib.Floor(o.R, 1)))
	binary.BigEndian.PutUint64(data[28:36], math.Float64bits(lib.Floor(o.Dir, 2)))
	binary.BigEndian.PutUint64(data[36:44], math.Float64bits(lib.Floor(o.Mh, 1)))
	binary.BigEndian.PutUint64(data[44:52], math.Float64bits(lib.Floor(o.H, 1)))
	binary.BigEndian.PutUint64(data[52:60], math.Float64bits(lib.Floor(o.Opacity, 2)))
	binary.BigEndian.PutUint32(data[60:64], uint32(o.Exp))
	if o.IsDead {
		data = append(data, 1)
	} else {
		data = append(data, 0)
	}
	data = append(data, byte(len(o.Team)))
	data = append(data, []byte(o.Team)...)
	data = append(data, byte(len(o.Type)))
	data = append(data, []byte(o.Type)...)
	if !o.IsShowName {
		o.Name = ""
	}
	data = append(data, byte(len(o.Name)))
	data = append(data, []byte(o.Name)...)

	return data
}

func (o *Object) SetController(p *Player) {
	o.Controller = p
}

func NewObject(value map[string]interface{}, t func(*Object), c func(*Object, *Object), k func(*Object, *Object), d func(*Object, *Object)) *Object {
	m := map[string]interface{}{
		"id":           objID,
		"type":         "squ",
		"team":         "ffa",
		"x":            -999999,
		"y":            -999999,
		"r":            lib.Grid,
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
		"spawnTime":    lib.Now(),
		"hitTime":      lib.Now(),
		"deadTime":     -1,
		"isBorder":     true,
		"isOwnCol":     true,
		"isDead":       false,
		"isCollision":  false,
		"isShowHealth": true,
		"isShowName":   false,
	}
	objID++

	for key, val := range value {
		m[key] = val
	}

	jsonString, _ := json.Marshal(m)

	s := Object{}
	json.Unmarshal(jsonString, &s)
	s.Tick = t
	s.Collision = c
	s.KillEvent = k
	s.DeadEvent = d
	return &s
}
