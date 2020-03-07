package main

import (
	"log"
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
		for _, obj := range objects {
			obj.X += obj.Dx
			obj.Y += obj.Dy
			if f, ok := obj.Event["Tick"].(func(obj *lib.Object)); ok {
				f(obj)
			}
		}
	}
}

func sendUpdates(ticker time.Ticker) {
	for range ticker.C {

	}
}
