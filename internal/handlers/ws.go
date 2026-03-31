package handlers

import (
	"net/http"
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan []byte)            // Broadcast channel
var mutex = &sync.Mutex{}                    // Protect clients map

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	Subprotocols:    []string{"chat"},
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "<http://localhost:3000>"
	},
}

type wsHandler struct{}

func NewWSHandler() RouteHandler {
	return &wsHandler{}
}

func (h *wsHandler) Register(app *fiber.App) {
	app.Get("/ws", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		mutex.Lock()
		clients[conn] = true
		mutex.Unlock()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				mutex.Lock()
				delete(clients, conn)
				mutex.Unlock()
				break
			}
			broadcast <- message
		}
	}))
}


func HandleWSMessages() {
	for {
		// Grab the next message from the broadcast channel
		message := <-broadcast

		// Send the message to all connected clients
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
