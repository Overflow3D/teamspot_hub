package hub

import "github.com/gorilla/websocket"

type WsConnections struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// NewConnection adds new ws
func NewConnection(hub *Hub, conn *websocket.Conn) *WsConnections {
	return &WsConnections{
		hub:  hub,
		conn: conn,
		send: make(chan []byte),
	}
}
