package main

import (
	//"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"app/event"
	"app/lib"

	"github.com/gorilla/websocket"
	//sq "github.com/mattn/go-sqlite3"
)

//var sockets = map[string]socketio.Conn{}
var objects []*lib.Object
var users = make(map[string]*lib.Player)

//var server, _ = socketio.NewServer(nil)
var quadtree lib.Quadtree
var scoreboard lib.Scoreboard
var setting lib.Setting = lib.ReadSetting()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}

	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))

		var objmap map[string]interface{}
		_ = json.Unmarshal(p, &objmap)
		event := objmap["event"].(string)

		switch event {
		case "init":
		case "input":

		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func serverWs(hub *event.Hub, w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &event.Client{Hub: hub, Conn: connection, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.ReadPump()
	go client.CatchMessage()
}

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
	if envPort != "" {
		port = envPort
	}

	lib.ShapeCount += setting.MaxShape

	hub := event.NewHub()
	go hub.Run()

	http.Handle("/", http.FileServer(http.Dir("./dist")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(hub, w, r)
	})

	// moveLoopTicker := time.NewTicker(time.Second / 60)
	// sendUpdatesTicker := time.NewTicker(time.Second / 30)

	// defer moveLoopTicker.Stop()
	// defer sendUpdatesTicker.Stop()

	// go moveloop(*moveLoopTicker)
	// go sendUpdates(*sendUpdatesTicker)

	log.Println("INFO > Server is Running Port" + port)
	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}

// func moveloop(ticker time.Ticker) {
// 	for range ticker.C {
// 		for ; lib.ShapeCount > 0; lib.ShapeCount-- {
// 			objects = append(objects, lib.NewObject(map[string]interface{}{
// 				"type":   "Square",
// 				"name":   "Square",
// 				"team":   "shape",
// 				"x":      lib.RandomRange(-setting.MapSize.X, setting.MapSize.X),
// 				"y":      lib.RandomRange(-setting.MapSize.Y, setting.MapSize.Y),
// 				"dir":    lib.RandomRange(-math.Pi, math.Pi),
// 				"stance": 0.2,
// 				"exp":    10,
// 			}, nil, nil, lib.DefaultCollision, nil, func(obj *lib.Object, killer *lib.Object) {
// 				lib.ShapeCount++
// 			}))
// 		}
// 		scoreboard = lib.Scoreboard{}

// 		for _, u := range users {
// 			u.PlayerSet()

// 			if u.ControlObject != nil {
// 				scoreboard.Push(lib.Score{
// 					Name:  u.ControlObject.Name,
// 					Type:  u.ControlObject.Type,
// 					Score: u.ControlObject.Exp,
// 				})
// 			}
// 		}

// 		quadtree = *lib.NewQuadtree(-setting.MapSize.X-lib.Grid*4, -setting.MapSize.Y-lib.Grid*4, setting.MapSize.X*2+lib.Grid*8, setting.MapSize.Y*2+lib.Grid*8, 1)

// 		for i := 0; i < len(objects); i++ {
// 			obj := objects[i]

// 			if obj.Tick != nil {
// 				obj.Tick(obj)
// 			}

// 			obj.ObjectTick()

// 			if obj.IsBorder { // 화면 밖으로 벗어나는가?
// 				if obj.X > setting.MapSize.X+lib.Grid*4 {
// 					obj.X = setting.MapSize.X + lib.Grid*4
// 				}
// 				if obj.X < -setting.MapSize.X-lib.Grid*4 {
// 					obj.X = -setting.MapSize.X - lib.Grid*4
// 				}
// 				if obj.Y > setting.MapSize.Y+lib.Grid*4 {
// 					obj.Y = setting.MapSize.Y + lib.Grid*4
// 				}
// 				if obj.Y < -setting.MapSize.Y-lib.Grid*4 {
// 					obj.Y = -setting.MapSize.Y - lib.Grid*4
// 				}
// 			}

// 			objList := quadtree.Retrieve(obj)

// 			for _, obj2 := range objList {
// 				if math.Sqrt((obj.X-obj2.X)*(obj.X-obj2.X)+(obj.Y-obj2.Y)*(obj.Y-obj2.Y)) < obj.R+obj2.R {
// 					obj.Collision(obj, obj2)
// 					obj2.Collision(obj2, obj)
// 				}
// 			}

// 			quadtree.Insert(obj)

// 			if obj.IsDead {
// 				if obj.DeadTime == -1 {
// 					obj.DeadTime = 700
// 				} else if obj.DeadTime <= 0 {
// 					if obj.DeadEvent != nil {
// 						obj.DeadEvent(obj, obj.HitObject)
// 					}
// 					objects = append(objects[:i], objects[i+1:]...)
// 					/*objects[i] = objects[len(objects)-1]
// 					objects = objects[:len(objects)-1]*/
// 					i--
// 				} else {
// 					obj.DeadTime = math.Max(obj.DeadTime-1000./60., 0.)
// 				}
// 			} else {
// 				obj.DeadTime = -1
// 			}
// 		}
// 	}
// }

// func sendUpdates(ticker time.Ticker) {
// 	for range ticker.C {
// 		var test lib.Quadtree = *lib.NewQuadtree(-setting.MapSize.X-lib.Grid*4, -setting.MapSize.Y-lib.Grid*4, setting.MapSize.X*2+lib.Grid*8, setting.MapSize.Y*2+lib.Grid*8, 1)
// 		for _, obj := range objects {
// 			test.Insert(obj)
// 		}
// 		for _, u := range users {
// 			objList := test.Retrieve(lib.Area{
// 				X: u.Camera.Pos.X - 1280/u.Camera.Z,
// 				Y: u.Camera.Pos.Y - 720/u.Camera.Z,
// 				W: 1280 / u.Camera.Z * 2,
// 				H: 720 / u.Camera.Z * 2,
// 			})
// 			//log.Println(objList)
// 			var visibleObject = []map[string]interface{}{}

// 			for _, obj := range objList {
// 				if obj.X < u.Camera.Pos.X+1280/u.Camera.Z+obj.R &&
// 					obj.X > u.Camera.Pos.X-1280/u.Camera.Z-obj.R &&
// 					obj.Y < u.Camera.Pos.Y+720/u.Camera.Z+obj.R &&
// 					obj.Y > u.Camera.Pos.Y-720/u.Camera.Z-obj.R && obj.Opacity > 0 {
// 					visibleObject = append(visibleObject, obj.SocketObj())
// 				}
// 			}

// 			//st := time.Now()

// 			/*
// 				if s, ok := sockets[u.ID].(socketio.Conn); ok {
// 					s.Emit("objectList", visibleObject)
// 					//log.Println(time.Since(st))
// 					if obj := u.ControlObject; obj == nil {
// 						s.Emit("playerSet", map[string]interface{}{
// 							"id":          "",
// 							"level":       1,
// 							"isCanRotate": u.IsCanDir,
// 							"stat":        0,
// 							"stats":       nil,
// 							"maxstats":    nil,
// 						}, u.Camera)
// 					} else {
// 						s.Emit("playerSet", map[string]interface{}{
// 							"id":          s.ID(),
// 							"level":       obj.Level,
// 							"isCanRotate": u.IsCanDir,
// 							"stat":        u.Stat,
// 							"stats":       obj.Stats,
// 							"maxstats":    obj.MaxStats,
// 						}, u.Camera)
// 					}
// 				}
// 			*/
// 		}
// 		/*
// 			server.BroadcastToRoom("", "game", "scoreboard", scoreboard)
// 			server.BroadcastToRoom("", "game", "area", []lib.Area{
// 				lib.Area{
// 					X: -setting.MapSize.X,
// 					Y: -setting.MapSize.Y,
// 					W: setting.MapSize.X * 2,
// 					H: setting.MapSize.Y * 2,
// 				},
// 			})
// 		*/
// 	}
// }
