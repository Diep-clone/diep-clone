package event

import (
	"bytes"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

//Sockets is
var Sockets map[string]*Client = make(map[string]*Client)

// const (
// 	maxMessageSize = 512
// 	writeWait      = time.Second
// 	timeWait       = time.Second
// 	timePeriod     = (timeWait * 9) / 10
// )

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

//Client is
type Client struct {
	Hub *Hub

	Conn *websocket.Conn

	ID string
}

var clientID int = 1

//NewClient is
func NewClient(h *Hub, c *websocket.Conn) *Client {
	client := Client{
		Hub:  h,
		Conn: c,
		ID:   strconv.Itoa(clientID),
	}

	clientID++
	return &client
}

//ReadPump is
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println(err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		Event(c, message)
	}

}

// WritePump is
// func (c *Client) WritePump() {
// 	ticker := time.NewTicker(timePeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.Conn.Close()
// 	}()

// 	for {
// 		select {
// 		case message, ok := <-c.Send:
// 			// c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			// if !ok {
// 				// c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				// return
// 			// }

// 			w, err := c.Conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				return
// 			}
// 			w.Write(message)

// 			n := len(c.Send)
// 			for i := 0; i < n; i++ {
// 				w.Write(newline)
// 				w.Write(<-c.Send)
// 			}

// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 		case <-ticker.C:
// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
