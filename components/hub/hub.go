package hub

import "fmt"

type Hub struct {
	connections map[*WsConnections]bool
	broadcast   chan []byte
	register    chan *WsConnections
	unregister  chan *WsConnections
}

func New() *Hub {
	return &Hub{
		connections: make(map[*WsConnections]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *WsConnections),
		unregister:  make(chan *WsConnections),
	}
}

// Run starts up hub
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.connections[conn] = true
			fmt.Println("Connected")
		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				close(conn.send)
			}
		case message := <-h.broadcast:
			for conn := range h.connections {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(h.connections, conn)
				}
			}
		}
	}
}
