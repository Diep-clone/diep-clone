package main

import (
	//"database/sql"

	"math"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"app/lib"
	"app/obj"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	//sq "github.com/mattn/go-sqlite3"
)

var quadtree obj.Quadtree
var setting lib.Setting = lib.ReadSetting()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serverWs(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithError(err).Error("WebSocket Connecting Error")
		return
	}

	client := obj.NewPlayer(connection)
	client.Register()

	go client.ReadPump()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	/*database, _ := sql.Open("sqlite3","./bogo.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id TEXT, nickname TEXT, password TEXT")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO people (nickname) VALUES (?)")
	statement.Exec("")
	*/
	rand.Seed(time.Now().UnixNano())

	port := "80"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	obj.ShapeCount = setting.MaxShape

	http.Handle("/", http.FileServer(http.Dir("./dist")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(w, r)
	})

	moveLoopTicker := time.NewTicker(time.Second / 60)
	sendUpdatesTicker := time.NewTicker(time.Second / 30)

	defer moveLoopTicker.Stop()
	defer sendUpdatesTicker.Stop()

	go moveloop(*moveLoopTicker)
	go sendUpdates(*sendUpdatesTicker)

	log.Info("Server is Running Port ", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

var su = new(sync.Mutex)
var isSend = false

func moveloop(ticker time.Ticker) {
	for range ticker.C {
		su.Lock()
		for ; obj.ShapeCount > 0; obj.ShapeCount-- {
			obj.Objects = append(obj.Objects, obj.NewObject(map[string]interface{}{
				"type":   "Square",
				"name":   "Square",
				"team":   "shape",
				"x":      lib.RandomRange(-lib.GameSetting.MapSize.X, lib.GameSetting.MapSize.X),
				"y":      lib.RandomRange(-lib.GameSetting.MapSize.Y, lib.GameSetting.MapSize.Y),
				"dir":    lib.RandomRange(-math.Pi, math.Pi),
				"stance": 0.2,
				"exp":    10,
			}, nil, nil, obj.DefaultCollision, nil, func(o *obj.Object, killer *obj.Object) {
				obj.ShapeCount++
			}))
		}
		scoreboard := lib.Scoreboard{}

		for _, u := range obj.Users {
			u.PlayerSet()

			if u.ControlObject != nil {
				scoreboard.Push(lib.Score{
					Name:  u.ControlObject.Name,
					Type:  u.ControlObject.Type,
					Score: u.ControlObject.Exp,
				})
			}
		}

		quadtree = *obj.NewQuadtree(-lib.GameSetting.MapSize.X-lib.Grid*4, -lib.GameSetting.MapSize.Y-lib.Grid*4, lib.GameSetting.MapSize.X*2+lib.Grid*8, lib.GameSetting.MapSize.Y*2+lib.Grid*8, 1)

		for i := 0; i < len(obj.Objects); i++ {
			o := obj.Objects[i]

			if o.Tick != nil {
				o.Tick(o)
			}

			o.ObjectTick()

			objList := quadtree.Retrieve(o)

			for _, obj2 := range objList {
				if math.Sqrt((o.X-obj2.X)*(o.X-obj2.X)+(o.Y-obj2.Y)*(o.Y-obj2.Y)) < o.R+obj2.R {
					o.Collision(o, obj2)
					obj2.Collision(obj2, o)
				}
			}

			quadtree.Insert(o)

			if o.IsDead {
				if o.DeadTime == -1 {
					o.DeadTime = 400
				} else if o.DeadTime <= 0 {
					if o.DeadEvent != nil {
						o.DeadEvent(o, o.HitObject)
					}
					obj.Objects = append(obj.Objects[:i], obj.Objects[i+1:]...)
					i--
				} else {
					o.DeadTime = math.Max(o.DeadTime-1000./60., 0.)
				}
			} else {
				o.DeadTime = -1
			}
		}

		su.Unlock()
		runtime.Gosched()
	}
}

var im = 50

func sendUpdates(ticker time.Ticker) {
	for range ticker.C {
		su.Lock()
		st := time.Now()

		for _, u := range obj.Users {
			u.CameraSet()
			var objList []*obj.Object = quadtree.Retrieve(obj.Area{
				X: u.Camera.Pos.X - 1280/u.Camera.Z,
				Y: u.Camera.Pos.Y - 720/u.Camera.Z,
				W: 1280 / u.Camera.Z * 2,
				H: 720 / u.Camera.Z * 2,
			})

			var visibleObject []map[string]interface{} = make([]map[string]interface{}, 0, len(objList)) // 미리 할당

			for _, o := range objList {
				if o.X < u.Camera.Pos.X+1280/u.Camera.Z+o.R &&
					o.X > u.Camera.Pos.X-1280/u.Camera.Z-o.R &&
					o.Y < u.Camera.Pos.Y+720/u.Camera.Z+o.R &&
					o.Y > u.Camera.Pos.Y-720/u.Camera.Z-o.R && o.Opacity > 0 {
					visibleObject = append(visibleObject, o.SocketObj())
				}
			}

			if u.Conn != nil {
				go SendUpdate(u, visibleObject)
			}
		}

		if im == 0 {
			log.Println(time.Since(st))
			im = 50
		}
		im--
		su.Unlock()
	}
}

func SendUpdate(u *obj.Player, visibleObject []map[string]interface{}) {
	if o := u.ControlObject; o == nil {
		u.Send(map[string]interface{}{
			"event":  "playerSet",
			"data":   false,
			"camera": u.Camera,
		})
	} else {
		u.Send(map[string]interface{}{
			"event":       "playerSet",
			"data":        true,
			"id":          u.ID,
			"level":       o.Level,
			"isCanRotate": u.IsCanDir,
			"stat":        u.Stat,
			"stats":       o.Stats,
			"maxstats":    o.MaxStats,
			"camera":      u.Camera,
		})
	}

	u.Send(map[string]interface{}{"event": "objectList", "data": visibleObject})
	u.Send(map[string]interface{}{
		"event": "area",
		"data": []obj.Area{
			{
				X: -setting.MapSize.X,
				Y: -setting.MapSize.Y,
				W: setting.MapSize.X * 2,
				H: setting.MapSize.Y * 2,
			},
		},
	})
}
