package obj

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"
	"sync"
	"unicode/utf8"

	"app/lib"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
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

var Users map[int]*Player = make(map[int]*Player)

// Player is
type Player struct {
	ID            int
	Conn          *websocket.Conn
	Mu            sync.Mutex
	Camera        Camera
	Mx            float64
	My            float64
	Mr            bool // right Mouse
	Ml            bool // left Mouse
	Stat          int
	MoveDir       float64
	IsMove        bool
	IsCanDir      bool
	StartTime     int64
	ControlObject *Object
	KillObject    *Object
}

var playerID int = 1

// NewPlayer is make player
func NewPlayer(c *websocket.Conn) *Player {
	p := Player{}
	p.ID = playerID
	playerID++
	p.Conn = c
	p.StartTime = lib.Now()
	p.IsCanDir = true

	return &p
}

// Player's mouse point set
func (p *Player) SetMousePoint(x float64, y float64) {
	p.Mx = x
	p.My = y
}

//
func (p *Player) CameraSet() {
	if obj := p.ControlObject; obj != nil {
		p.Camera = Camera{
			Pos: Pos{X: lib.Floor(obj.X, 2), Y: lib.Floor(obj.Y, 2)},
			Z:   16. / 9. * math.Pow(0.995, float64(obj.Level)-1) / float64(obj.Sight),
		}
	}
}

//
func (p *Player) PlayerSet() {
	if obj := p.ControlObject; obj != nil {
		if p.IsMove {
			obj.Dx += math.Cos(p.MoveDir) * obj.Speed
			obj.Dy += math.Sin(p.MoveDir) * obj.Speed
		}
		if p.IsCanDir {
			obj.Dir = math.Atan2(p.My, p.Mx)
		}
		if obj.DeadTime == 0 {
			p.Send([]byte{2})
			p.ControlObject = nil
		}
	}
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

//
func (p *Player) Send(v []byte) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	return p.Conn.WriteMessage(websocket.BinaryMessage, v)
}

// Register == connect
func (p *Player) Register() {
	log.WithField("id", p.ID).Info("User Connected")
	p.Send(AreaSocket())
}

//ReadPump is
func (p *Player) ReadPump() {
	defer func() {
		p.UnRegister()
		p.Conn.Close()
	}()

	for {
		_, message, err := p.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithFields(log.Fields{
					"message": message,
					"error":   err,
				}).Error("WebSocket Close Error")
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		Event(p, message)
	}

}

// Event Catch Message
func Event(p *Player, message []byte) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	var socketType uint8 = uint8(message[0])

	switch socketType {
	case 0:
		if _, ok := Users[p.ID]; ok {
			log.WithField("id", p.ID).Warn("Prevent Login")
			return
		}
		Users[p.ID] = p
	case 1:
		if u, ok := Users[p.ID]; ok {
			dir := math.Float32frombits(binary.BigEndian.Uint32(message[1:5]))
			u.IsMove = (dir != 9.)
			u.MoveDir = float64(dir)

			var x, y float32 = math.Float32frombits(binary.BigEndian.Uint32(message[5:9])), math.Float32frombits(binary.BigEndian.Uint32(message[9:13]))

			if obj := u.ControlObject; obj != nil {
				u.SetMousePoint(float64(x)-obj.X, float64(y)-obj.Y)
			}

			var value uint8 = uint8(message[13])

			u.Ml = (value % 2) == 1
			value /= 2
			u.Mr = (value % 2) == 1
			value /= 2
			if value%2 == 1 {
				if obj := u.ControlObject; obj != nil {
					obj.H = 0
				}
			}
		}
	case 2:
		var test []byte = message[1:16]

		var name string = ""
		for len(test) > 0 {
			r, size := utf8.DecodeRune(test)
			if string(r) != "\x00" {
				name += string(r)
			}
			test = test[size:]
		}
		log.Println(name)

		p.ControlObject = NewObject(map[string]interface{}{
			"type":       "Necromanser",
			"team":       strconv.Itoa(p.ID),
			"name":       name,
			"x":          lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
			"y":          lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
			"h":          50,
			"mh":         50,
			"damage":     20,
			"level":      45,
			"stats":      [8]int{0, 0, 0, 7, 7, 7, 7, 5},
			"maxStats":   [8]int{7, 7, 7, 7, 7, 7, 7, 7},
			"sight":      1.11,
			"isShowName": true,
		}, []Gun{*NewGun(map[string]interface{}{
			"owner": nil,
			"limit": 0,
		}, nil, Object{})}, TankTick, DefaultCollision, NecroKillEvent, nil)
		p.ControlObject.SetController(p)
		Objects = append(Objects, p.ControlObject)
	default:
		log.WithFields(log.Fields{
			"id": p.ID,
		}).Error("socketType Error")
		return
	}
}

// UnRegister == close
func (p *Player) UnRegister() {
	if _, ok := Users[p.ID]; !ok {
		log.WithField("id", p.ID).Warn("Prevent Disconnect")
		return
	}

	if p.ControlObject != nil {
		p.ControlObject.Controller = nil
	}

	delete(Users, p.ID)

	log.WithField("id", p.ID).Info("User Disconnect")
}
