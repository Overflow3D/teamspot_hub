package hub

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

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

func (ws *WsConnections) ListnerAndReader() {
	ws.hub.register <- ws
	go ws.readMessages()
	go ws.writeMessages()
}

func (ws *WsConnections) readMessages() {
	defer func() {
		ws.hub.unregister <- ws
		ws.conn.Close()
	}()
	ws.conn.SetReadLimit(maxMessageSize)
	ws.conn.SetReadDeadline(time.Now().Add(pongWait))
	ws.conn.SetPongHandler(func(string) error { ws.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		ws.hub.broadcast <- message
	}
}

func (ws *WsConnections) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.conn.Close()
	}()
	for {
		select {
		case message, ok := <-ws.send:
			ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				ws.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := ws.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(ws.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-ws.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
