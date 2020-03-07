package lib

import "time"

// Camera is
type Camera struct {
	Pos Pos
	Z   float64
}

// Player is
type Player struct {
	Id            string
	Mx            float64
	My            float64
	Camera        Camera
	Keys          map[string]interface{}
	StartTime     int64
	ControlObject *Object
}

// NewPlayer is make player
func NewPlayer(id string) *Player {
	p := Player{}
	p.Id = id
	p.Mx = 0
	p.My = 0
	p.Camera = Camera{
		Pos: Pos{0, 0},
		Z:   1,
	}
	p.Keys = map[string]interface{}{}
	p.StartTime = time.Now().Unix()

	return &p
}

// Player's mouse point set
func (p *Player) SetMousePoint(x float64, y float64) {
	p.Mx = x
	p.My = y
}

// Player's key set
func (p *Player) SetKey(key string, value interface{}) {
	p.Keys[key] = value

}
