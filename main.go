package main

import (
	//"database/sql"

	"encoding/binary"
	"math"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
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
	scoreBoardTicker := time.NewTicker(time.Second / 100)

	defer moveLoopTicker.Stop()
	defer sendUpdatesTicker.Stop()
	defer scoreBoardTicker.Stop()

	go moveloop(*moveLoopTicker)
	go sendUpdates(*sendUpdatesTicker)
	go scoreBoard(*scoreBoardTicker)

	log.Info("Server is Running Port ", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

var su = new(sync.Mutex)

func moveloop(ticker time.Ticker) {
	for range ticker.C {
		su.Lock()

		obj.AddShape()

		for _, u := range obj.Users {
			u.PlayerSet()
		}

		quadtree = *obj.NewQuadtree(-lib.GameSetting.MapSize.X-lib.Grid*4, -lib.GameSetting.MapSize.Y-lib.Grid*4, lib.GameSetting.MapSize.X*2+lib.Grid*8, lib.GameSetting.MapSize.Y*2+lib.Grid*8, 1)

		for i := 0; i < len(obj.Objects); i++ {
			o := obj.Objects[i]

			if o.Tick != nil {
				o.Tick(o)
			}

			o.ObjectTick()

			if !o.IsDead {
				objList := quadtree.Retrieve(obj.Area{
					X: o.X - o.R,
					Y: o.Y - o.R,
					W: o.R * 2,
					H: o.R * 2,
				})

				for _, obj2 := range objList {
					if !obj2.IsDead && (obj2.Owner != o.Owner || obj2.IsOwnCol && o.IsOwnCol) && o != obj2.Owner && obj2 != o.Owner {
						if math.Sqrt((o.X-obj2.X)*(o.X-obj2.X)+(o.Y-obj2.Y)*(o.Y-obj2.Y)) < o.R+obj2.R {
							o.Collision(o, obj2)
							obj2.Collision(obj2, o)
						}
					}
				}
			}

			quadtree.Insert(o)

			if o.IsDead {
				if o.DeadTime == -1 {
					o.DeadTime = 300
				} else if o.DeadTime <= 0 {
					if o.DeadEvent != nil {
						o.DeadEvent(o, o.HitObject)
					}
					obj.Objects = append(obj.Objects[:i], obj.Objects[i+1:]...)
					i--
				} else {
					o.DeadTime = math.Max(o.DeadTime-1000./60., 0.)
				}
			}
		}

		su.Unlock()
	}
}

func sendUpdates(ticker time.Ticker) {
	for range ticker.C {
		su.Lock()
		for _, u := range obj.Users {
			u.CameraSet()

			var objList []*obj.Object = quadtree.Retrieve(obj.Area{
				X: u.Camera.Pos.X - 960/u.Camera.Z,
				Y: u.Camera.Pos.Y - 540/u.Camera.Z,
				W: 960 / u.Camera.Z * 2,
				H: 540 / u.Camera.Z * 2,
			})

			var sendData []byte = make([]byte, 32)

			binary.BigEndian.PutUint64(sendData[1:9], math.Float64bits(u.Camera.Pos.X))
			binary.BigEndian.PutUint64(sendData[9:17], math.Float64bits(u.Camera.Pos.Y))
			binary.BigEndian.PutUint64(sendData[17:25], math.Float64bits(u.Camera.Z))

			sendData[25] = 0
			if o := u.ControlObject; o != nil {
				sendData[25] = 1
				binary.BigEndian.PutUint32(sendData[26:30], uint32(o.ID))
				binary.BigEndian.PutUint16(sendData[30:32], uint16(o.Level))
				sendData = append(sendData, byte(u.Stat))
				for i := 0; i < 8; i++ {
					sendData = append(sendData, byte(o.Stats[i]))
				}
				for i := 0; i < 8; i++ {
					sendData = append(sendData, byte(o.MaxStats[i]))
				}
				sendData = append(sendData, byte(len(o.Team)))
				sendData = append(sendData, []byte(o.Team)...)
			}

			for _, o := range objList {
				if o.X < u.Camera.Pos.X+960/u.Camera.Z+o.R &&
					o.X > u.Camera.Pos.X-960/u.Camera.Z-o.R &&
					o.Y < u.Camera.Pos.Y+540/u.Camera.Z+o.R &&
					o.Y > u.Camera.Pos.Y-540/u.Camera.Z-o.R && o.Opacity > 0 {
					sendData = append(sendData, o.ObjBinaryData()...)
				}
			}

			if u.Conn != nil {
				err := u.Send(sendData)
				if err != nil {
					log.Println("Send Error")
					return
				}
			}
		}
		su.Unlock()
	}
}

func scoreBoard(ticker time.Ticker) {
	for range ticker.C {
		su.Lock()

		var scoreBoard lib.Scoreboard

		for _, u := range obj.Users {
			if o := u.ControlObject; o != nil {
				scoreBoard.Push(lib.Score{
					Team:  o.Team,
					Name:  o.Name,
					Type:  o.Type,
					Score: o.Exp,
				})
			}
		}

		var sendData []byte = make([]byte, 5)

		sendData[0] = 3

		for _, v := range scoreBoard {
			var data []byte = make([]byte, 4)
			binary.BigEndian.PutUint32(data[0:4], uint32(v.Score))

			sendData = append(sendData, data...)
			sendData = append(sendData, byte(len(v.Team)))
			sendData = append(sendData, []byte(v.Team)...)
			sendData = append(sendData, byte(len(v.Type)))
			sendData = append(sendData, []byte(v.Type)...)
			sendData = append(sendData, byte(len(v.Name)))
			sendData = append(sendData, []byte(v.Name)...)
		}

		for _, u := range obj.Users {
			if u.Conn != nil {
				err := u.Send(sendData)
				if err != nil {
					log.Println("send Error")
					return
				}
			}
		}
		su.Unlock()
	}
}
