package main

import (
	"encoding/binary"
	"math"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	"app/lib"
	"app/obj"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

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

// TODO : Achievement
// TODO : GameSpeed
// TODO : GameMode

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
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

var count = 0
var count2 = 0

func moveloop(ticker time.Ticker) { // manages the motion of all objects.
	for range ticker.C {
		obj.ObjMutex.Lock()

		obj.AddShape() // Spawn Shapes

		for _, u := range obj.Users { // Users setting
			u.PlayerSet()
			u.CameraSet()
		}

		obj.Qt.Clear() // Clear Quadtree

		for _, o := range obj.Objects { // Collision Quadtree
			if !o.IsDead {
				obj.Qt.Insert(o)
			}
		}

		for i, o := range obj.Objects {
			o.Index = i
			if o.Tick != nil {
				o.Tick(o)
			}

			o.ObjectTick()
		}

		for _, o := range obj.Objects {
			for _, e := range obj.Qt.Retrieve(obj.Area{
				X: o.X - o.R,
				Y: o.Y - o.R,
				W: o.R * 2,
				H: o.R * 2,
			}) {
				if o != e && !o.IsDead && !e.IsDead && !e.IsCollision && (e.Owner != o.Owner || e.IsOwnCol && o.IsOwnCol) && o != e.Owner && e != o.Owner && !lib.Contains(e.HitObjects, o) {
					if (o.X-e.X)*(o.X-e.X)+(o.Y-e.Y)*(o.Y-e.Y) < (o.R+e.R)*(o.R+e.R) {
						if o.Collision != nil {
							o.Collision(o, e)
						}
						if e.Collision != nil {
							e.Collision(e, o)
						}
					}
				}
			}
		}

		for i := 0; i < len(obj.Objects); i++ {
			var o *obj.Object = obj.Objects[i]
			o.Index = i
			if len(o.HitObjects) > 0 {
				o.LastHitObject = o.HitObjects[0]
			}
			if o.IsDead {
				if o.DeadTime == -1 {
					o.DeadTime = 300
				} else if o.DeadTime <= 0 {
					if o.DeadEvent != nil {
						o.DeadEvent(o, o.LastHitObject)
					}
					obj.ObjIDList = append(obj.ObjIDList, o.ID)
					obj.Objects = append(obj.Objects[:o.Index], obj.Objects[o.Index+1:]...)
					i--
				} else {
					o.Opacity = math.Max(o.Opacity-0.1, 0) // dead effect
					o.R += o.R * 0.02
					o.DeadTime = math.Max(o.DeadTime-1000./60., 0.)
				}
			}
			o.HitObjects = nil
		}

		obj.ObjMutex.Unlock()
	}
}

func sendUpdates(ticker time.Ticker) {
	for range ticker.C {
		obj.ObjMutex.Lock()

		for _, u := range obj.Users {
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

			for _, o := range obj.Objects {
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
					log.Println("Send Error to sendUpdates")
					continue
				}
			}
		}

		obj.ObjMutex.Unlock()
	}
}

func scoreBoard(ticker time.Ticker) {
	for range ticker.C {
		obj.ObjMutex.Lock()

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

		var sendData []byte = make([]byte, 1)

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
					log.Println("send Error to scoreboard")
					continue
				}
			}
		}
		obj.ObjMutex.Unlock()
	}
}
