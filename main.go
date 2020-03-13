package main

import (
	//"database/sql"
	"log"
	"math"
	"net/http"
	"time"

	"app/lib"

	socketio "github.com/googollee/go-socket.io"
	//sq "github.com/mattn/go-sqlite3"
)

var setting lib.Setting = lib.ReadSetting()

var sockets = make(map[string]interface{})
var objects []*lib.Object
var users = make(map[string]*lib.Player)

var server, _ = socketio.NewServer(nil)
var quadtree lib.Quadtree

var scoreboard lib.Scoreboard

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU()) 모든 CPU를 사용하게 해주는 코드
	/*database, _ := sql.Open("sqlite3","./bogo.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id TEXT, nickname TEXT, password TEXT")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO people (nickname) VALUES (?)")
	statement.Exec("")
	*/
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
			3,
			"ffa",
			name,
			lib.RandomRange(-setting.MapSize.X, setting.MapSize.Y),
			lib.RandomRange(-setting.MapSize.Y, setting.MapSize.X),
			lib.Grid,
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
		users[s.ID()].ControlObject.Basic()
		objects = append(objects, users[s.ID()].ControlObject)

		log.Println("INFO > " + s.ID() + " Login")
	})

	server.OnEvent("/", "keyboard", func(s socketio.Conn, key string, enable bool, value interface{}) {
		if u, ok := users[s.ID()]; ok {
			if enable {
				u.SetKey(key, value)
				if obj := u.ControlObject; obj != nil {
					if f, ok := u.ControlObject.Event["KeyDown"].(map[string]func())[key]; ok {
						f()
					}
				}
			} else {
				u.DisableKey(key)
				if obj := u.ControlObject; obj != nil {
					if f, ok := u.ControlObject.Event["KeyUp"].(map[string]func())[key]; ok {
						f()
					}
				}
			}
		}
	})

	server.OnEvent("/", "stat", func(s socketio.Conn, n int) {
		if u, ok := users[s.ID()]; ok {
			if obj := u.ControlObject; obj != nil {
				if obj.Variable["Stat"].(int) > 0 && obj.Variable["Stats"].([]int)[n] < obj.Variable["MaxStats"].([]int)[n] {
					obj.Variable["Stats"].([]int)[n]++
					if v, ok := obj.Variable["Stat"].(int); ok {
						v--
					}
				}
			}
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

		scoreboard = lib.Scoreboard{}

		for _, u := range users {

			if obj := u.ControlObject; obj != nil {
				for key, _ := range obj.Event {
					if f, ok := u.ControlObject.Event["KeyPress"].(map[string]func())[key]; ok {
						f()
					}
				}
			}

			scoreboard.Push(lib.Score{
				Name:  u.ControlObject.Name,
				Type:  u.ControlObject.Type,
				Color: u.ControlObject.Color,
				Score: u.ControlObject.Exp,
			})
		}

		quadtree.Clear()

		for _, obj := range objects {
			obj.C.Pos.X += obj.Dx
			obj.C.Pos.Y += obj.Dy

			obj.Dx *= 0.97
			obj.Dy *= 0.97

			if obj.IsDead {
				continue
			}

			if obj.H <= 0 {
				obj.IsDead = true
			}

			if obj.IsBorder { // 화면 밖으로 벗어나는가?
				if obj.C.Pos.X > setting.MapSize.X+lib.Grid*4 {
					obj.C.Pos.X = setting.MapSize.X + lib.Grid*4
				}
				if obj.C.Pos.X < -setting.MapSize.X-lib.Grid*4 {
					obj.C.Pos.X = -setting.MapSize.X - lib.Grid*4
				}
				if obj.C.Pos.Y > setting.MapSize.Y+lib.Grid*4 {
					obj.C.Pos.Y = setting.MapSize.Y + lib.Grid*4
				}
				if obj.C.Pos.Y < -setting.MapSize.Y-lib.Grid*4 {
					obj.C.Pos.Y = -setting.MapSize.Y - lib.Grid*4
				}
			}

			if time.Now().Unix()-30000 > obj.HitTime {
				obj.H += obj.Mh / 60 / 10
			}

			imMh := obj.Mh

			if f, ok := obj.Event["Tick"].(func(obj *lib.Object)); !obj.IsDead && ok {
				f(obj)
			}

			objList := quadtree.Retrieve(obj)

			for _, obj2 := range objList {
				if math.Sqrt((obj.C.Pos.X-obj2.C.Pos.X)*(obj.C.Pos.X-obj2.C.Pos.X)+(obj.C.Pos.Y-obj2.C.Pos.Y)*(obj.C.Pos.Y-obj2.C.Pos.Y)) < obj.C.R+obj2.C.R {
					if lib.CollisionEvent(obj, obj2) {
						server.BroadcastToRoom("/", "/", "collision", obj.ID)
						server.BroadcastToRoom("/", "/", "collision", obj2.ID)
					}
				}
			}

			quadtree.Insert(obj)

			if obj.H > obj.Mh {
				obj.H = obj.Mh
			}

			if obj.Mh != imMh {
				obj.H *= obj.Mh / imMh
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
			var visibleObject = []map[string]interface{}{}
			for _, obj := range objList {
				if obj.C.Pos.X < u.Camera.Pos.X+1280/u.Camera.Z+obj.C.R &&
					obj.C.Pos.X > u.Camera.Pos.X+1280/u.Camera.Z-obj.C.R &&
					obj.C.Pos.Y < u.Camera.Pos.Y+720/u.Camera.Z+obj.C.R &&
					obj.C.Pos.Y > u.Camera.Pos.Y+720/u.Camera.Z-obj.C.R && obj.Opacity > 0 {
					visibleObject = append(visibleObject, obj.SocketObj())
				}
			}

			if s, ok := sockets[u.ID].(socketio.Conn); ok {
				s.Emit("objectList", visibleObject)
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
