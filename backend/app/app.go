package app

import (
	"localVercel/internal/queue"
	"localVercel/internal/runtime"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type App struct {
	startedAt time.Time
	upgrader  websocket.Upgrader
	Queue     queue.Queue
	Runtime   *runtime.Manager
}

func New(q queue.Queue) *App {
	return &App{
		Queue:     q,
		Runtime:   runtime.NewManager(),
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
