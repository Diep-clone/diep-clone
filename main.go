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

var sockets = map[string]socketio.Conn{}
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
	quadtree = *lib.NewQuadtree(-setting.MapSize.X-lib.Grid*4, -setting.MapSize.Y-lib.Grid*4, setting.MapSize.X*2+lib.Grid*8, setting.MapSize.Y*2+lib.Grid*8, 1)

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
			"controller": *users[s.ID()],
			"type":       "Necromanser",
			"name":       name,
			"x":          lib.RandomRange(-setting.MapSize.X, setting.MapSize.X),
			"y":          lib.RandomRange(-setting.MapSize.Y, setting.MapSize.Y),
			"h":          50,
			"mh":         50,
			"damage":     20,
			"level":      45,
			"stats":      [8]int{0, 0, 0, 7, 7, 7, 5, 7},
			"maxStats":   [8]int{7, 7, 7, 7, 7, 7, 7, 7},
			"sight":      1.11,
			/*"guns": []lib.Gun{*lib.NewGun(map[string]interface{}{
				"owner": users[s.ID()].ControlObject,
				"limit": 0,
			})},*/
		}, func(obj *lib.Object) {
			lib.TankTick(obj)
		}, lib.DefaultCollision, nil, /*func(a *lib.Object, b *lib.Object) {
				if b.Type == "Square" && a.Guns[0].Limit < 22+2*a.Stats[6] {
					b.Type = "NecroSquare"
					b.Team = a.Team
					b.Exp = 0
					b.Owner = a
					b.Mh = (8 + 6*float64(a.Stats[4])) * 0.5
					b.H = b.Mh
					b.Damage = (7 + 3*float64(a.Stats[5])) * 1.68
					b.Speed = (0.056 + 0.01*float64(a.Stats[3])) * 0.5
					b.IsDead = false
					b.KillEvent = a.KillEvent
					b.DeadEvent = func(obj *lib.Object, killer *lib.Object) {
						obj.Owner.Guns[0].Limit--
						lib.ShapeCount++
					}
					b.Tick = func(obj *lib.Object) {
						p := obj.Owner.Controller
						if p == nil {
							return
						}
						if p.Mr {
							obj.Dir = math.Atan2(obj.Y-(obj.Owner.Y+p.My), obj.X-(obj.Owner.X+p.Mx))
							obj.Dx += math.Cos(obj.Dir) * obj.Speed
							obj.Dy += math.Sin(obj.Dir) * obj.Speed
						} else if p.Ml {
							obj.Dir = math.Atan2((obj.Owner.Y+p.My)-obj.Y, (obj.Owner.X+p.Mx)-obj.X)
							obj.Dx += math.Cos(obj.Dir) * obj.Speed
							obj.Dy += math.Sin(obj.Dir) * obj.Speed
						} else if obj.Owner.DeadTime != 0 {

						} else {
							obj.H = 0
						}
					}
				}
			}*/nil)
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

	server.OnEvent("/", "mouseright", func(s socketio.Conn, b bool) {
		if u, ok := users[s.ID()]; ok {
			u.Mr = b
		}
	})

	server.OnEvent("/", "mouseleft", func(s socketio.Conn, b bool) {
		if u, ok := users[s.ID()]; ok {
			u.Ml = b
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
		for ; lib.ShapeCount > 0; lib.ShapeCount-- {
			objects = append(objects, lib.NewObject(map[string]interface{}{
				"type": "Square",
				"name": "Square",
				"team": "shape",
				"x":    lib.RandomRange(-setting.MapSize.X, setting.MapSize.X),
				"y":    lib.RandomRange(-setting.MapSize.Y, setting.MapSize.Y),
				"dir":  lib.RandomRange(-math.Pi, math.Pi),
			}, nil, lib.DefaultCollision, nil, func(obj *lib.Object, killer *lib.Object) {
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

		quadtree.Clear()

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
					obj.DeadTime = 300
				} else if obj.DeadTime < 0 {
					if obj.DeadEvent != nil {
						obj.DeadEvent(obj, obj.HitObject)
					}
					objects[i] = objects[len(objects)-1]
					objects = objects[:len(objects)-1]
					i--
				} else {
					obj.DeadTime = math.Min(obj.DeadTime-1000./60., 0.)
				}
			}
		}
	}
}

func sendUpdates(ticker time.Ticker) {
	for range ticker.C {
		for _, u := range users {
			/*objList := quadtree.Retrieve(lib.Area{
				X:  u.Camera.Pos.X - 1280/u.Camera.Z,
				Y:  u.Camera.Pos.Y - 720/u.Camera.Z,
				X2: 1280/u.Camera.Z * 2,
				Y2: 720/u.Camera.Z * 2,
			})
			log.Println(objList)*/
			var visibleObject = []map[string]interface{}{}

			for _, obj := range objects {
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
						"id":          -1,
						"level":       1,
						"isCanRotate": u.IsCanDir,
						"stat":        0,
						"stats":       nil,
						"maxstats":    nil,
					}, u.Camera)
				} else {
					s.Emit("playerSet", map[string]interface{}{
						"id":          obj.ID,
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
