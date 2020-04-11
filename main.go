package main

import (
	//"database/sql"
	"os"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"app/lib"

	socketio "github.com/googollee/go-socket.io"
	//sq "github.com/mattn/go-sqlite3"
)

var sockets = map[string]socketio.Conn{}
var objects []*lib.Object
var users = make(map[string]*lib.Player)

var server, _ = socketio.NewServer(nil)
var quadtree lib.Quadtree

var scoreboard lib.Scoreboard

var setting lib.Setting = lib.ReadSetting()

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU()) 모든 CPU를 사용하게 해주는 코드
	/*database, _ := sql.Open("sqlite3","./bogo.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id TEXT, nickname TEXT, password TEXT")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO people (nickname) VALUES (?)")
	statement.Exec("")
	*/
	rand.Seed(time.Now().UnixNano())
	
	port := "3000"
	envPort := os.Getenv("PORT")
	if envPort != nil {
	    port = envPort
	}

	lib.ShapeCount += setting.MaxShape

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("INFO > " + s.ID() + " Connected")
		s.Join("game")

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

		sockets[s.ID()] = s
		users[s.ID()] = lib.NewPlayer(s.ID())
		users[s.ID()].ControlObject = lib.NewObject(map[string]interface{}{
			"type":     "Necromanser",
			"team":     s.ID(),
			"name":     name,
			"x":        lib.RandomRange(-setting.MapSize.X, setting.MapSize.X),
			"y":        lib.RandomRange(-setting.MapSize.Y, setting.MapSize.Y),
			"h":        50,
			"mh":       50,
			"damage":   20,
			"level":    45,
			"stats":    [8]int{0, 0, 0, 7, 7, 7, 7, 5},
			"maxStats": [8]int{7, 7, 7, 7, 7, 7, 7, 7},
			"sight":    1.11,
		}, []lib.Gun{*lib.NewGun(map[string]interface{}{
			"owner": users[s.ID()].ControlObject,
			"limit": 0,
		})}, lib.TankTick, lib.DefaultCollision, lib.NecroKillEvent, nil)
		users[s.ID()].ControlObject.SetController(users[s.ID()])
		objects = append(objects, users[s.ID()].ControlObject)

		log.Println("INFO > " + s.ID() + " Login")
	})

	server.OnEvent("/", "moveVector", func(s socketio.Conn, value float64) {
		if u, ok := users[s.ID()]; ok {
			u.IsMove = value != 9
			u.MoveDir = value
		}
	})

	server.OnEvent("/", "stat", func(s socketio.Conn, n int) {

	})

	server.OnEvent("/", "mousemove", func(s socketio.Conn, x float64, y float64) {
		if u, ok := users[s.ID()]; ok {
			if obj := u.ControlObject; obj != nil {
				u.SetMousePoint(x-obj.X, y-obj.Y)
			} else {
				u.SetMousePoint(x, y)
			}
		}
	})

	server.OnEvent("/", "mouseleft", func(s socketio.Conn, b bool) {
		if u, ok := users[s.ID()]; ok {
			u.Ml = b
		}
	})

	server.OnEvent("/", "mouseright", func(s socketio.Conn, b bool) {
		if u, ok := users[s.ID()]; ok {
			u.Mr = b
		}
	})

	server.OnEvent("/", "kill", func(s socketio.Conn, b bool) {
		if u, ok := users[s.ID()]; ok {
			if u.ControlObject != nil {
				u.ControlObject.H = 0
			}
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

		users[s.ID()].ControlObject.Controller = nil

		delete(sockets, s.ID())
		delete(users, s.ID())

		log.Println("INFO > " + s.ID() + " disconnected")
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./dist")))

	moveLoopTicker := time.NewTicker(time.Second / 60)
	sendUpdatesTicker := time.NewTicker(time.Second / 30)

	defer moveLoopTicker.Stop()
	defer sendUpdatesTicker.Stop()

	go moveloop(*moveLoopTicker)
	go sendUpdates(*sendUpdatesTicker)

	log.Println("INFO > Server is Running Port" + port)
	log.Fatal(http.ListenAndServe("localhost:" + port, nil))
}

func moveloop(ticker time.Ticker) {
	for range ticker.C {
		for ; lib.ShapeCount > 0; lib.ShapeCount-- {
			objects = append(objects, lib.NewObject(map[string]interface{}{
				"type":   "Square",
				"name":   "Square",
				"team":   "shape",
				"x":      lib.RandomRange(-setting.MapSize.X, setting.MapSize.X),
				"y":      lib.RandomRange(-setting.MapSize.Y, setting.MapSize.Y),
				"dir":    lib.RandomRange(-math.Pi, math.Pi),
				"stance": 0.2,
				"exp":    10,
			}, nil, nil, lib.DefaultCollision, nil, func(obj *lib.Object, killer *lib.Object) {
				lib.ShapeCount++
			}))
		}
		scoreboard = lib.Scoreboard{}

		for _, u := range users {
			u.PlayerSet()

			if u.ControlObject != nil {
				scoreboard.Push(lib.Score{
					Name:  u.ControlObject.Name,
					Type:  u.ControlObject.Type,
					Score: u.ControlObject.Exp,
				})
			}
		}

		quadtree = *lib.NewQuadtree(-setting.MapSize.X-lib.Grid*4, -setting.MapSize.Y-lib.Grid*4, setting.MapSize.X*2+lib.Grid*8, setting.MapSize.Y*2+lib.Grid*8, 1)

		for i := 0; i < len(objects); i++ {
			obj := objects[i]

			if obj.Tick != nil {
				obj.Tick(obj)
			}

			obj.ObjectTick()

			if obj.IsBorder { // 화면 밖으로 벗어나는가?
				if obj.X > setting.MapSize.X+lib.Grid*4 {
					obj.X = setting.MapSize.X + lib.Grid*4
				}
				if obj.X < -setting.MapSize.X-lib.Grid*4 {
					obj.X = -setting.MapSize.X - lib.Grid*4
				}
				if obj.Y > setting.MapSize.Y+lib.Grid*4 {
					obj.Y = setting.MapSize.Y + lib.Grid*4
				}
				if obj.Y < -setting.MapSize.Y-lib.Grid*4 {
					obj.Y = -setting.MapSize.Y - lib.Grid*4
				}
			}

			objList := quadtree.Retrieve(obj)

			for _, obj2 := range objList {
				if math.Sqrt((obj.X-obj2.X)*(obj.X-obj2.X)+(obj.Y-obj2.Y)*(obj.Y-obj2.Y)) < obj.R+obj2.R {
					obj.Collision(obj, obj2)
					obj2.Collision(obj2, obj)
				}
			}

			quadtree.Insert(obj)

			if obj.IsDead {
				if obj.DeadTime == -1 {
					obj.DeadTime = 700
				} else if obj.DeadTime <= 0 {
					if obj.DeadEvent != nil {
						obj.DeadEvent(obj, obj.HitObject)
					}
					objects = append(objects[:i], objects[i+1:]...)
					/*objects[i] = objects[len(objects)-1]
					objects = objects[:len(objects)-1]*/
					i--
				} else {
					obj.DeadTime = math.Max(obj.DeadTime-1000./60., 0.)
				}
			} else {
				obj.DeadTime = -1
			}
		}
	}
}

func sendUpdates(ticker time.Ticker) {
	for range ticker.C {
		var test lib.Quadtree = *lib.NewQuadtree(-setting.MapSize.X-lib.Grid*4, -setting.MapSize.Y-lib.Grid*4, setting.MapSize.X*2+lib.Grid*8, setting.MapSize.Y*2+lib.Grid*8, 1)
		for _, obj := range objects {
			test.Insert(obj)
		}
		for _, u := range users {
			objList := test.Retrieve(lib.Area{
				X: u.Camera.Pos.X - 1280/u.Camera.Z,
				Y: u.Camera.Pos.Y - 720/u.Camera.Z,
				W: 1280 / u.Camera.Z * 2,
				H: 720 / u.Camera.Z * 2,
			})
			//log.Println(objList)
			var visibleObject = []map[string]interface{}{}

			for _, obj := range objList {
				if obj.X < u.Camera.Pos.X+1280/u.Camera.Z+obj.R &&
					obj.X > u.Camera.Pos.X-1280/u.Camera.Z-obj.R &&
					obj.Y < u.Camera.Pos.Y+720/u.Camera.Z+obj.R &&
					obj.Y > u.Camera.Pos.Y-720/u.Camera.Z-obj.R && obj.Opacity > 0 {
					visibleObject = append(visibleObject, obj.SocketObj())
				}
			}

			//st := time.Now()

			if s, ok := sockets[u.ID].(socketio.Conn); ok {
				s.Emit("objectList", visibleObject)
				//log.Println(time.Since(st))
				if obj := u.ControlObject; obj == nil {
					s.Emit("playerSet", map[string]interface{}{
						"id":          "",
						"level":       1,
						"isCanRotate": u.IsCanDir,
						"stat":        0,
						"stats":       nil,
						"maxstats":    nil,
					}, u.Camera)
				} else {
					s.Emit("playerSet", map[string]interface{}{
						"id":          s.ID(),
						"level":       obj.Level,
						"isCanRotate": u.IsCanDir,
						"stat":        u.Stat,
						"stats":       obj.Stats,
						"maxstats":    obj.MaxStats,
					}, u.Camera)
				}
			}
		}
		server.BroadcastToRoom("", "game", "scoreboard", scoreboard)
		server.BroadcastToRoom("", "game", "area", []lib.Area{
			lib.Area{
				X: -setting.MapSize.X,
				Y: -setting.MapSize.Y,
				W: setting.MapSize.X * 2,
				H: setting.MapSize.Y * 2,
			},
		})
	}
}
