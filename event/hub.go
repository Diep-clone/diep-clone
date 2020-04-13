package event

// Hub 가 클라이언트의 집합성 유지 및 관리
type Hub struct {
	Clients map[*Client]bool

	Broadcast chan []byte

	// 생성
	Register chan *Client

	// 삭제
	Unregister chan *Client
}

// NewHub make new hub
func NewHub() *Hub {

	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

//Run is Run
func (h *Hub) Run() {

	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			Register(client)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				UnRegister(client)

				delete(h.Clients, client)
				close(client.Send)
			}
		}
	}
}
