package obj

import (
	"bytes"
	"encoding/json"
	"math"
	"strconv"
	"sync"

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

var Users map[string]*Player = make(map[string]*Player)

// Player is
type Player struct {
	ID            string
	Conn          *websocket.Conn
	mu            sync.Mutex
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

var playerID int = 1

// NewPlayer is make player
func NewPlayer(c *websocket.Conn) *Player {
	p := Player{}
	p.ID = strconv.Itoa(playerID)
	playerID++
	p.Conn = c
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

//
func (p *Player) CameraSet() {
	if obj := p.ControlObject; obj != nil {
		p.Camera.Pos = Pos{X: obj.X, Y: obj.Y}
		p.Camera.Z = 16. / 9. * math.Pow(0.995, float64(obj.Level)-1) / float64(obj.Sight)
	} else {
		p.Camera.Pos = Pos{X: 0, Y: 0}
		p.Camera.Z = 1
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
	}
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

//
func (p *Player) Send(v interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Conn.WriteJSON(v)
}

// Register == connect
func (p *Player) Register() {
	log.WithField("id", p.ID).Info("User Connected")
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
	//fmt.Println(string(message))

	var objmap map[string]interface{}

	if err := json.Unmarshal(message, &objmap); err != nil {
		log.WithFields(log.Fields{
			"id":    p.ID,
			"error": err,
		}).Error("JSON Unmarshal Error")
		return
	}

	event, ok := objmap["event"].(string)
	if !ok {
		return
	}

	//fmt.Println(objmap)

	//fmt.Println(event)

	switch event {
	case "init":
		if _, ok := Users[p.ID]; ok {
			log.WithField("id", p.ID).Warn("Prevent Login")
			return
		}

		name, ok := objmap["data"].(string)
		if !ok {
			return
		}

		if len(name) > 15 {
			name = ""
		}

		Users[p.ID] = p
		/*
			var gunList []obj.Gun
			for i := -math.Pi / 2; i <= math.Pi; i += math.Pi / 2 {
				gunList = append(gunList, *obj.NewGun(map[string]interface{}{
					"dir": i,
				}, obj.Users[c.ID].ControlObject, *obj.NewObject(map[string]interface{}{
					"speed":  1,
					"h":      1,
					"damage": 0.65,
					"r":      1,
					"stance": 0.1,
				}, []obj.Gun{*obj.NewGun(map[string]interface{}{
					"delaytime": 3.,
				}, nil, obj.Object{})}, obj.DefaultBulletTick, obj.DefaultCollision, nil, nil)))
			}
			for i := -math.Pi / 4 * 3; i <= math.Pi; i += math.Pi / 2 {
				gunList = append(gunList, *obj.NewGun(map[string]interface{}{
					"dir":      i,
					"waittime": 0.5,
				}, obj.Users[c.ID].ControlObject, *obj.NewObject(map[string]interface{}{
					"speed":  1,
					"h":      1,
					"damage": 0.65,
					"r":      1,
					"stance": 0.1,
				}, []obj.Gun{*obj.NewGun(map[string]interface{}{
					"delaytime": 3.,
				}, nil, obj.Object{})}, obj.DefaultBulletTick, obj.DefaultCollision, nil, nil)))
			}
			obj.Users[c.ID].ControlObject = obj.NewObject(map[string]interface{}{
				"type":     "OctoTank",
				"team":     string(c.ID),
				"name":     name,
				"x":        lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
				"y":        lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
				"h":        50,
				"mh":       50,
				"damage":   20,
				"level":    45,
				"stats":    [8]int{0, 0, 0, 7, 7, 7, 7, 5},
				"maxStats": [8]int{7, 7, 7, 7, 7, 7, 7, 7},
			}, gunList, obj.TankTick, obj.DefaultCollision, nil, nil)
		*/
		Users[p.ID].ControlObject = NewObject(map[string]interface{}{
			"type":     "Necromanser",
			"team":     string(p.ID),
			"name":     name,
			"x":        lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
			"y":        lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
			"h":        50,
			"mh":       50,
			"damage":   20,
			"level":    45,
			"stats":    [8]int{0, 0, 0, 7, 7, 7, 7, 5},
			"maxStats": [8]int{7, 7, 7, 7, 7, 7, 7, 7},
			"sight":    1.11,
		}, []Gun{*NewGun(map[string]interface{}{
			"owner": nil,
			"limit": 0,
		}, nil, Object{})}, TankTick, DefaultCollision, NecroKillEvent, nil)
		Users[p.ID].ControlObject.SetController(Users[p.ID])
		Objects = append(Objects, Users[p.ID].ControlObject)

	case "input":
		if m, ok := objmap["data"].(map[string]interface{}); ok {
			if t, ok := m["type"].(string); ok {
				if u, ok := Users[p.ID]; ok {
					switch t {
					case "moveVector":
						if value, ok := m["value"].(float64); ok {
							u.IsMove = value != 9
							u.MoveDir = value
						}
					case "mouseMove":
						pos, ok := m["value"].(map[string]interface{})
						if !ok {
							return
						}
						if obj := u.ControlObject; obj != nil {
							u.SetMousePoint(pos["x"].(float64)-obj.X, pos["y"].(float64)-obj.Y)
						} else {
							u.SetMousePoint(pos["x"].(float64), pos["y"].(float64))
						}
					case "mouseLeft":
						if value, ok := m["value"].(bool); ok {
							u.Ml = value
						} else {
							u.Ml = false
						}
					case "mouseRight":
						if value, ok := m["value"].(bool); ok {
							u.Mr = value
						} else {
							u.Mr = false
						}
					case "suicide":
						if obj := u.ControlObject; obj != nil {
							obj.H = 0
						}
					}
				}
			}
		}
	}
}

// UnRegister == close
func (p *Player) UnRegister() {
	if _, ok := Users[p.ID]; !ok {
		log.WithField("id", p.ID).Warn("Prevent Disconnect")
		return
	}

	Users[p.ID].ControlObject.Controller = nil

	delete(Users, p.ID)

	log.WithField("id", p.ID).Info("User Disconnect")
}
