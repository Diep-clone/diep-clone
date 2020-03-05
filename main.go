package main

import (
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

var sockets map[string]interface{}

// 그렇군요 어제 그 강의 싹다 봤긴 했는데 어...음음 그럼 대충
func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU()) 모든 CPU를 사용하게 해주는 코드
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	// 뭐 해요 setinterval이요? 그건 setinterval을 고루틴으로 만들수 있겠죠 뭐 ㅁ
	server.OnConnect("/", func(s socketio.Conn) error {
		log.Println("INFO > " + s.ID() + " connected")
		/*p := lib.newPlayer(s.ID())
		println(p)*/
		return nil
	})

	server.OnEvent("/", "ping!", func(s socketio.Conn, data interface{}) {
		s.Emit("pong!", data)
	})

	server.OnDisconnect("/", func(s socketio.Conn, _ string) {
		log.Println("INFO > " + s.ID() + " disconnected")
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("INFO > Server is Running Port 3000")
	log.Fatal(http.ListenAndServe("localhost:3000", nil))

	go moveloop()
	go sendUpdates()
}

func moveloop() {
	for {

	}
}

func sendUpdates() {
	for {

	}
}
