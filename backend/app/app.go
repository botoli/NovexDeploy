package app

import (
	"localVercel/internal/queue"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type App struct {
	startedAt time.Time
	upgrader  websocket.Upgrader
	Queue     queue.Queue
}

func New(q queue.Queue) *App {
	return &App{
		Queue:     q,
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