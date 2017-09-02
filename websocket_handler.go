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

func wshandler(h *hub.Hub, w http.ResponseWriter, r *http.Request) {
	_, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
}
