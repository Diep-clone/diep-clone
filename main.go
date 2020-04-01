package main

import (
	//"database/sql"
	"log"
	"net/http"
	"time"

	"app/lib"

	socketio "github.com/googollee/go-socket.io"
	//sq "github.com/mattn/go-sqlite3"
)

var setting lib.Setting = lib.ReadSetting()

var sockets = map[string]socketio.Conn{}
var tanks []lib.Tank
var bullets []lib.Bullet
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

		sockets[s.ID()] = s
		users[s.ID()] = lib.NewPlayer(s.ID())
		tanks = append(tanks, *lib.NewTank(map[string]interface{}{
			"Controller": users[s.ID()],
			"Type":       "Circle",
			"Color":      3,
			"Team":       "ffa",
			"Name":       name,
			"X":          lib.RandomRange(-setting.MapSize.X, setting.MapSize.Y),
			"Y":          lib.RandomRange(-setting.MapSize.Y, setting.MapSize.X),
			"R":          lib.Grid,
		}))

		log.Println("INFO > " + s.ID() + " Login")
	})

	server.OnEvent("/", "moveVector", func(s socketio.Conn, value float64) {
		users[s.ID()].MoveDir = value
	})

	server.OnEvent("/", "stat", func(s socketio.Conn, n int) {

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

		/*
			for _, u := range users {

				scoreboard.Push(lib.Score{
					Name:  u.ControlObject.Name,
					Type:  u.ControlObject.Type,
					Color: u.ControlObject.Color,
					Score: u.ControlObject.Exp,
				})
			}
		*/

		quadtree.Clear()
		/*
			for _, obj := range objects {
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
						// Collision Event
					}
				}

				quadtree.Insert(obj)
			}
		*/
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
			log.Println(objList)
			var visibleObject = []map[string]interface{}{}
			for _, obj := range tanks {
				if obj.X < u.Camera.Pos.X+1280/u.Camera.Z+obj.R &&
					obj.X > u.Camera.Pos.X-1280/u.Camera.Z-obj.R &&
					obj.Y < u.Camera.Pos.Y+720/u.Camera.Z+obj.R &&
					obj.Y > u.Camera.Pos.Y-720/u.Camera.Z-obj.R && obj.Opacity > 0 {
					visibleObject = append(visibleObject, obj.SocketObj())
				}
			}

			if s, ok := sockets[u.ID].(socketio.Conn); ok {
				s.Emit("objectList", visibleObject)
				s.Emit("playerSet", map[string]interface{}{
					"id":          u.ControlObjectID,
					"level":       1,      // u.ControlObject.Variable["Level"],
					"sight":       16 / 9, // u.ControlObject.Variable["Sight"],
					"isCanRotate": u.IsCanDir,
					"stat":        0,   // u.ControlObject.Variable["Stat"],
					"stats":       nil, // u.ControlObject.Variable["Stats"],
					"maxstats":    nil, // u.ControlObject.Variable["MaxStats"],
				}, u.Camera)
			}
		}
		server.BroadcastToRoom("/", "/", "scoreboard", scoreboard)
	}
}
