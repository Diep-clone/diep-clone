package main

import (
	"log"
	"math"
	"net/http"
	"time"

	"app/lib"

	socketio "github.com/googollee/go-socket.io"
)

var setting = lib.ReadSetting()

var sockets = make(map[string]interface{})
var objects []*lib.Object
var users = make(map[string]*lib.Player)

var server, _ = socketio.NewServer(nil)

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU()) 모든 CPU를 사용하게 해주는 코드
	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("INFO > " + s.ID() + " Connected")

		return nil
	})

	server.OnEvent("/", "ping!", func(s socketio.Conn, data interface{}) {
		s.Emit("pong!", data)
	})

	server.OnEvent("/", "login", func(s socketio.Conn, name string) {
		if _, ok := sockets[s.ID()]; ok {
			log.Println("INFO > Prevent " + s.ID() + " Login")
			return
		}
		if len(name) > 15 {
			name = ""
		}

		sockets[s.ID()] = nil
		users[s.ID()] = lib.NewPlayer(s.ID())
		users[s.ID()].ControlObject = lib.NewObject(
			users[s.ID()],
			[]interface{}{},
			[3]int{0, 0, 0},
			"ffa",
			name,
			lib.RandomRange(-setting.MapSize.X, setting.MapSize.Y),
			lib.RandomRange(-setting.MapSize.Y, setting.MapSize.X),
			lib.RandomRange(100, 100),
			lib.RandomRange(100, 100),
			lib.RandomRange(100, 100),
			lib.RandomRange(100, 100),
			lib.RandomRange(100, 100),
			lib.RandomRange(100, 100),
			map[string]interface{}{
				"Tick": func(obj *lib.Object) { // obj
					obj.Speed = (0.07 + (0.007 * obj.Variable["Stats"].([]float64)[7])) * math.Pow(0.985, obj.Variable["Level"].(float64)-1)
					obj.Damage = (20 + obj.Variable["Stats"].([]float64)[2]*4)
					obj.R = (13 * math.Pow(1.01, (obj.Variable["Level"].(float64)-1)))
					obj.Mh = (48 + obj.Variable["Level"].(float64)*2 + obj.Variable["Stats"].([]float64)[1]*20)
				},
				"KeyDown":   nil, // obj, keyType
				"KeyPress":  nil, // obj, keyType, time
				"KeyUp":     nil, // obj, keyType
				"GetBound":  nil, // obj, enemyObj
				"GetDamage": nil, // obj, enemyObj, damage
				"KillEvent": nil, // obj, enemyObj
				"DeadEvent": nil, // obj, enemyObj
			},
			map[string]interface{}{
				"Level":         1,
				"MaxStats":      []int{7, 7, 7, 7, 7, 7, 7, 7},
				"Stats":         []int{0, 0, 0, 0, 0, 0, 0, 0},
				"Stat":          0,
				"Sight":         1.,
				"InvisibleTime": 0.,
			},
			false,
			false,
		)
		objects = append(objects, users[s.ID()].ControlObject)

		log.Println("INFO > " + s.ID() + " Login")
	})

	server.OnEvent("/", "keyboard", func(s socketio.Conn, key string, time float64) {
		if u, ok := users[s.ID()]; ok {
			u.SetKey(key, time)
		}
	})

	server.OnEvent("/", "mousemove", func(s socketio.Conn, mouse struct {
		x float64
		y float64
	}) {
		if u, ok := users[s.ID()]; ok {
			u.SetMousePoint(mouse.x, mouse.y)
		}
	})

	server.OnError("/", func(s socketio.Conn, err error) {
		log.Println(err)
	})

	server.OnDisconnect("/", func(s socketio.Conn, _ string) {
		if _, ok := sockets[s.ID()]; !ok {
			log.Println("INFO > Prevent " + s.ID() + " Disconnect")
			return
		}

		delete(sockets, s.ID())
		delete(users, s.ID())

		log.Println("INFO > " + s.ID() + " disconnected")
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./dist")))

	moveLoopTicker := time.NewTicker(time.Second / 60)
	sendUpdatesTicker := time.NewTicker(time.Second / 40)

	defer moveLoopTicker.Stop()
	defer sendUpdatesTicker.Stop()

	go moveloop(*moveLoopTicker)
	go sendUpdates(*sendUpdatesTicker)

	log.Println("INFO > Server is Running Port 3000")
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}

func moveloop(ticker time.Ticker) {
	for range ticker.C {
		for _, user := range users {
			if user.ControlObject != nil {
				if user.IsMove {
					user.ControlObject.Dx += math.Cos(user.Keys["moveRotate"]) * user.ControlObject.Speed
					user.ControlObject.Dy += math.Sin(user.Keys["moveRotate"]) * user.ControlObject.Speed
				}

			}
		}
	}
}

func sendUpdates(ticker time.Ticker) {
	for range ticker.C {

	}
}
