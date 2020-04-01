package lib

import "time"

// Camera is
type Camera struct {
	Pos Pos
	Z   float64
}

// Player is
type Player struct {
	ID              string
	Mx              float64
	My              float64
	Stat            int
	Stats           [8]int
	Camera          Camera
	MoveDir         float64
	IsCanDir        bool
	StartTime       int64
	ControlObjectID int
}

// NewPlayer is make player
func NewPlayer(id string) *Player {
	p := Player{}
	p.ID = id
	p.Mx = 0
	p.My = 0
	p.Stat = 0
	p.Stats = [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	p.Camera = Camera{
		Pos: Pos{0, 0},
		Z:   1,
	}
	p.StartTime = time.Now().Unix()

	return &p
}

// Player's mouse point set
func (p *Player) SetMousePoint(x float64, y float64) {
	p.Mx = x
	p.My = y
}

func (p *Player) SetCamera() {
	/*if p.ControlObject != nil {
		obj := *p.ControlObject
		p.Camera.Pos = Pos{X: obj.X, Y: obj.Y}
		p.Camera.Z = 16 / 9 // * math.Pow(0.995, obj.Variable["Level"].(float64)) / obj.Variable["Sight"].(float64)
	}*/
}
