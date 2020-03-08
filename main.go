package main

import (
	"log"
	"math"
	"net/http"
	"time"

	"app/lib"

	socketio "github.com/googollee/go-socket.io"
)

var setting lib.Setting = lib.ReadSetting()

var sockets = make(map[string]interface{})
var objects []*lib.Object
var users = make(map[string]*lib.Player)

var server, _ = socketio.NewServer(nil)
var quadtree lib.Quadtree

var scoreboard []lib.Scoreboard

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU()) 모든 CPU를 사용하게 해주는 코드
	quadtree = *lib.NewQuadtree(-1000, -1000, 2000, 2000, 1)

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
			"",
			"ffa",
			name,
			lib.RandomRange(-setting.MapSize.X, setting.MapSize.Y),
			lib.RandomRange(-setting.MapSize.Y, setting.MapSize.X),
			13,
			48,
			20,
			0.07,
			1,
			1,
			nil,
			nil,
			false,
			false,
		)
		users[s.ID()].ControlObject.Tank()
		objects = append(objects, users[s.ID()].ControlObject)

		log.Println("INFO > " + s.ID() + " Login")
	})

	server.OnEvent("/", "keyboard", func(s socketio.Conn, key string, time interface{}) {
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

		scoreboard := []lib.Scoreboard{}
		quadtree.Clear()

		for _, obj := range objects {
			if obj.H <= 0 {
				obj.IsDead = true
			}

			obj.X += obj.Dx
			obj.Y += obj.Dy

			if obj.IsDead {
				continue
			}

			if _, ok := (*obj.Owner).(lib.Player); ok {
				scoreboard = append(scoreboard, lib.Scoreboard{
					Name:  obj.Name,
					Type:  obj.Type,
					Score: obj.Exp,
				})
			}

			if obj.IsBorder { // 화면 밖으로 벗어나는가?
				if obj.X > setting.MapSize.X+51.6 {
					obj.X = setting.MapSize.X + 51.6
				}
				if obj.X < -setting.MapSize.X-51.6 {
					obj.X = -setting.MapSize.X - 51.6
				}
				if obj.Y > setting.MapSize.Y+51.6 {
					obj.Y = setting.MapSize.Y + 51.6
				}
				if obj.Y < -setting.MapSize.Y-51.6 {
					obj.Y = -setting.MapSize.Y - 51.6
				}
			}

			if time.Now().Unix()-30000 > obj.HitTime {
				obj.H += obj.Mh / 60 / 10
			}

			if f, ok := obj.Event["Tick"].(func(obj *lib.Object)); !obj.IsDead && ok {
				f(obj)
			}

			objList := quadtree.Retrieve(obj)

			for _, obj2 := range objList {
				if math.Sqrt((obj.X-obj2.X)*(obj.X-obj2.X)+(obj.Y-obj2.Y)*(obj.Y-obj2.Y)) < obj.R+obj2.R {
					if lib.CollisionEvent(obj, obj2) {
						// socket send collision Event
					}
				}
			}

			quadtree.Insert(obj)

			if obj.H > obj.Mh {
				obj.H = obj.Mh
			}
		}
	}
}

func sendUpdates(ticker time.Ticker) {
	for range ticker.C {
		for _, u := range users {
			objList := quadtree.Retrieve(lib.Area{
				X:  u.Camera.Pos.X + 1280/u.Camera.Z,
				Y:  u.Camera.Pos.Y + 720/u.Camera.Z,
				X2: u.Camera.Pos.X - 1280/u.Camera.Z,
				Y2: u.Camera.Pos.Y - 720/u.Camera.Z,
			})
			var visibleObj []map[string]interface{}
			for _, obj := range objList {
				if obj.X < u.Camera.Pos.X+1280/u.Camera.Z+obj.R &&
					obj.X > u.Camera.Pos.X+1280/u.Camera.Z-obj.R &&
					obj.Y < u.Camera.Pos.Y+720/u.Camera.Z+obj.R &&
					obj.Y > u.Camera.Pos.Y+720/u.Camera.Z-obj.R && obj.Opacity > 0 {
					visibleObj = append(visibleObj, obj.SocketObj())
				}
			}

			if s, ok := sockets[u.ID].(socketio.Conn); ok {
				s.Emit("objectList", visibleObj)
				s.Emit("playerSet", map[string]interface{}{
					"level":    u.ControlObject.Variable["Level"],
					"camera":   u.Camera,
					"isRotate": u.ControlObject.IsCanDir,
					"stat":     u.ControlObject.Variable["Stat"],
					"stats":    u.ControlObject.Variable["Stats"],
					"maxstats": u.ControlObject.Variable["MaxStats"],
				})
			}
		}
		server.BroadcastToRoom("/", "/", "scoreboard", scoreboard)
	}
}
