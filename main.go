package main

import (
	"github.com/Overflow3D/teamspot_hub/components/hub"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	hub := hub.New()
	go hub.Run()
	r.GET("/ws", wsHandler(hub))
	r.Run(":1333")
}

func wsHandler(h *hub.Hub) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		handleWS(h, ctx.Writer, ctx.Request)
	}
}
