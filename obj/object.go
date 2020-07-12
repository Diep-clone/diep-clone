package obj

import (
	"sync"

	"app/lib"
	"encoding/binary"
	"encoding/json"
	"math"
)

var Objects []*Object

var ObjMutex = new(sync.Mutex)

// Object is
type Object struct {
	Controller *Player `json:"controller"`
	Owner      *Object `json:"owner"`
	ID         int     `json:"id"`
	Type       string  `json:"type"`
	Team       string  `json:"team"`
	Name       string  `json:"name"`
	Index      int

	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	R      float64 `json:"r"`
	Dx     float64 `json:"dx"`
	Dy     float64 `json:"dy"`
	Dir    float64 `json:"dir"`
	Level  int     `json:"level"`
	Exp    int     `json:"exp"`
	H      float64 `json:"h"`
	Mh     float64 `json:"mh"`
	Lh     float64 `json:"lh"` // last frame health
	Damage float64 `json:"damage"`
	GetDH  float64 `json:"getdh"` // get damage and health
	Speed  float64 `json:"speed"`
	/*
		바운드 값은 상대방을 밀어내는 정도,
		스탠스 값은 자신이 밀쳐지는 정도를 뜻합니다.
		바운드/ 스탠스 시스템은 하자가 있습니다.
		연구 결과 데미지 값이 바운드/ 스탠스 값에 영향을 주는 사실을 발견했습니다.
		따라서 바운드/ 스탠스 시스템은 데미지와 관련된 값으로 재설계가 필요합니다.
	*/
	Bounce  float64 `json:"bound"`
	Stance  float64 `json:"stance"`
	Opacity float64 `json:"opacity"`

	Sight     float64 `json:"sight"`
	Stats     [8]int  `json:"stats"`
	MaxStats  [8]int  `json:"maxStats"`
	Guns      []Gun   `json:"guns"`
	SpawnTime int64   `json:"spawnTime"`
	HitTime   int64   `json:"hitTime"`
	DeadTime  float64 `json:"deadTime"`

	IsBorder     bool `json:"isBorder"`
	IsOwnCol     bool `json:"isOwnCol"` // 부모 오브젝트와의 충돌감지가 되는가?
	IsDead       bool `json:"isDead"`
	IsCollision  bool `json:"isCollision"` // Is finished collision detect?
	IsShot       bool `json:"isShot"`
	IsTargeted   bool `json:"isTargeted"`
	IsShowHealth bool `json:"isShowHealth"`
	IsShowName   bool `json:"isShowName"`
	IsBack       bool // 드론이 주인에게 돌아가야 하는가?

	/*
		서브 오브젝트는 오브젝트에 붙어있는 오브젝트들입니다.
		서브 오브젝트들은 히트박스, 체력값을 가지고 있지 않습니다.
	*/
	SubObjects []*Object
	HitObject  *Object
	Target     *Object
	Tick       func(*Object)
	Collision  func(*Object, *Object)
	KillEvent  func(*Object, *Object)
	DeadEvent  func(*Object, *Object)
}

//
func (obj *Object) ObjectTick() {
	obj.X += obj.Dx //* lib.GameSetting.GameSpeed // 좌표값에 속도 적용
	obj.Y += obj.Dy //* lib.GameSetting.GameSpeed

	obj.Dx *= 0.97 //math.Pow(0.97, lib.GameSetting.GameSpeed) 감속 코드
	obj.Dy *= 0.97 //math.Pow(0.97, lib.GameSetting.GameSpeed)

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

	if obj.H <= 0 && obj.Mh != 0 {
		obj.IsDead = true
		obj.H = 0
	}
	for _, o := range obj.SubObjects {
		if o == nil {
			continue
		}
		o.Index = -1
		if o.Tick != nil {
			o.Tick(o)
		}

		o.ObjectTick()
	}

	if obj.IsDead {
		return
	}

	obj.Shot()

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

	a.Dx += math.Cos(dir) * math.Min(b.Bounce*a.Stance, 6)
	a.Dy += math.Sin(dir) * math.Min(b.Bounce*a.Stance, 6)

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

//
func (o Object) ObjBinaryData() []byte {
	var data []byte = make([]byte, 64)
	if !o.IsShowHealth {
		o.Mh = 1
		o.H = 1
	}
	if !o.IsShowName {
		o.Name = ""
		o.Exp = 0
	}
	binary.BigEndian.PutUint32(data[0:4], uint32(o.ID))
	binary.BigEndian.PutUint64(data[4:12], math.Float64bits(o.X))
	binary.BigEndian.PutUint64(data[12:20], math.Float64bits(o.Y))
	binary.BigEndian.PutUint64(data[20:28], math.Float64bits(lib.Floor(o.R, 1)))
	binary.BigEndian.PutUint64(data[28:36], math.Float64bits(lib.Floor(o.Dir, 2)))
	binary.BigEndian.PutUint64(data[36:44], math.Float64bits(lib.Floor(o.Mh, 1)))
	binary.BigEndian.PutUint64(data[44:52], math.Float64bits(lib.Floor(o.H, 1)))
	binary.BigEndian.PutUint64(data[52:60], math.Float64bits(lib.Floor(o.Opacity, 2)))
	binary.BigEndian.PutUint32(data[60:64], uint32(o.Exp))
	data = append(data, byte(len(o.Guns)))
	for i := 0; i < len(o.Guns); i++ {
		g := &o.Guns[i]
		if g.IsShot {
			data = append(data, byte(1))
			g.IsShot = false
		} else {
			data = append(data, byte(0))
		}
	}
	data = append(data, byte(len(o.SubObjects)))
	for _, obj := range o.SubObjects {
		if obj == nil {
			data = append(data, make([]byte, 4)...)
		} else {
			data = append(data, obj.ObjBinaryData()...)
		}
	}
	if o.HitTime+100-lib.Now() > 0 {
		data = append(data, byte(uint8(o.HitTime+100-lib.Now())))
	} else {
		data = append(data, byte(0))
	}
	if o.IsDead {
		data = append(data, 1)
	} else {
		data = append(data, 0)
	}
	data = append(data, byte(len(o.Team)))
	data = append(data, []byte(o.Team)...)
	data = append(data, byte(len(o.Type)))
	data = append(data, []byte(o.Type)...)
	data = append(data, byte(len(o.Name)))
	data = append(data, []byte(o.Name)...)

	return data
}

func (o *Object) SetController(p *Player) {
	o.Controller = p
}

var objID = 1
var ObjIDList []int // 썼던 오브젝트 id 리스트. 재활용함.

func NewObject(value map[string]interface{}, t func(*Object), c func(*Object, *Object), k func(*Object, *Object), d func(*Object, *Object)) *Object {
	ID := objID
	objID++
	if len(ObjIDList) > 0 {
		ID = ObjIDList[0]
		if len(ObjIDList) < 2 {
			ObjIDList = nil
		} else {
			ObjIDList = ObjIDList[1 : len(ObjIDList)-1]
		}
		objID--
	}
	m := map[string]interface{}{
		"id":           ID,
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
		"hitTime":      0,
		"deadTime":     -1,
		"isBorder":     true,
		"isOwnCol":     true,
		"isDead":       false,
		"isTargeted":   true,
		"isCollision":  true,
		"isShowHealth": true,
		"isShowName":   false,
	}

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
