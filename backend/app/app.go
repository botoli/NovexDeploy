package app

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type App struct {
	startedAt time.Time
	upgrader  websocket.Upgrader
}

func New() *App {
	return &App{
		startedAt: time.Now(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
	}
}