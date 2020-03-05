package main

import (
	"log"
	"net/http"
	"time"

	"app/lib"

	socketio "github.com/googollee/go-socket.io"
)

var sockets map[string]interface{}
var users map[string]interface{}

var server, _ = socketio.NewServer(nil)

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU()) 모든 CPU를 사용하게 해주는 코드

	server.OnConnect("/", func(s socketio.Conn) error {
		if _, ok := sockets[s.ID()]; !ok {
			log.Println("INFO > Prevent " + s.ID() + " Login")
			return nil
		}

		sockets[s.ID()] = nil
		log.Println("INFO > " + s.ID() + " connected")

		users[s.ID()] = lib.Newplayer(s.ID())
		return nil
	})

	server.OnEvent("/", "ping!", func(s socketio.Conn, data interface{}) {
		s.Emit("pong!", data)
	})

	server.OnDisconnect("/", func(s socketio.Conn, _ string) {
		if _, ok := sockets[s.ID()]; !ok {
			log.Println("INFO > Prevent " + s.ID() + " Disconnect")
			return
		}

		delete(sockets, s.ID())
		log.Println("INFO > " + s.ID() + " disconnected")
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./static")))

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
		for u := range users {

		}
	}
}

func sendUpdates(ticker time.Ticker) {
	for range ticker.C {

	}
}
