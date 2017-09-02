package main

import (
	"net/http"

	"github.com/Overflow3D/teamspot_hub/components/hub"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleWS(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	wsConn := hub.NewConnection(h, conn)
	wsConn.ListnerAndReader()

}
