package obj

import (
	"math"

	"app/lib"
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

// Camera is
type Camera struct {
	Pos Pos
	Z   float64
}

var Users map[string]*Player = make(map[string]*Player)

// Player is
type Player struct {
	ID            string
	Mx            float64
	My            float64
	Mr            bool // right Mouse
	Ml            bool // left Mouse
	Stat          int
	Camera        Camera
	MoveDir       float64
	IsMove        bool
	IsCanDir      bool
	StartTime     int64
	ControlObject *Object
}

// NewPlayer is make player
func NewPlayer(id string) *Player {
	p := Player{}
	p.ID = id
	p.Camera = Camera{
		Pos: Pos{0, 0},
		Z:   1.,
	}
	p.StartTime = lib.Now()
	p.IsCanDir = true

	return &p
}

// Player's mouse point set
func (p *Player) SetMousePoint(x float64, y float64) {
	p.Mx = x
	p.My = y
}

func (p *Player) PlayerSet() {
	if obj := p.ControlObject; obj != nil {
		p.Camera.Pos = Pos{X: obj.X, Y: obj.Y}
		p.Camera.Z = 16. / 9. * math.Pow(0.995, float64(obj.Level)-1) / float64(obj.Sight)

		if p.IsMove {
			obj.Dx += math.Cos(p.MoveDir) * obj.Speed
			obj.Dy += math.Sin(p.MoveDir) * obj.Speed
		}
		if p.IsCanDir {
			obj.Dir = math.Atan2(p.My, p.Mx)
		}
	} else {
		p.Camera.Pos = Pos{X: 0, Y: 0}
		p.Camera.Z = 1
	}
}
