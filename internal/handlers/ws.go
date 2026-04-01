package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"putra4648/my-chat-app/internal/models"
	"putra4648/my-chat-app/internal/services"
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gofiber/fiber/v3/middleware/session"
	uuid "github.com/gofrs/uuid/v5"
	"github.com/gorilla/websocket"
)

var userConnections = make(map[string]*websocket.Conn) // Map of userID -> connection
var connMutex = &sync.Mutex{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Improved security would check origin here
	},
}

type wsHandler struct {
	chatService *services.ChatService
}

func NewWSHandler(chatService *services.ChatService) RouteHandler {
	return &wsHandler{chatService: chatService}
}

type WsMessage struct {
	Type           string    `json:"type"` // "text", "join", etc.
	ConversationID uuid.UUID `json:"conversation_id"`
	ReceiverID     uuid.UUID `json:"receiver_id"`
	Content        string    `json:"content"`
	SenderID       uuid.UUID `json:"sender_id"`
}

func (h *wsHandler) Register(app *fiber.App) {
	app.Get("/ws", func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		userID := sess.Get("user_id")
		if userID == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		userIDStr := userID.(string)

		handler := adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			defer conn.Close()

			connMutex.Lock()
			userConnections[userIDStr] = conn
			connMutex.Unlock()

			defer func() {
				connMutex.Lock()
				delete(userConnections, userIDStr)
				connMutex.Unlock()
			}()

			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					break
				}

				var wsMsg WsMessage
				if err := json.Unmarshal(message, &wsMsg); err != nil {
					continue
				}

				// Handle different message types
				if wsMsg.Type == "text" {
					// Save to DB
					msg := &models.Message{
						ConversationID: wsMsg.ConversationID,
						SenderID:       &wsMsg.SenderID,
						Content:        wsMsg.Content,
					}
					err := h.chatService.SaveMessage(context.Background(), msg)
					if err != nil {
						continue
					}

					// Forward to receiver if online
					receiverIDStr := wsMsg.ReceiverID.String()
					connMutex.Lock()
					receiverConn, online := userConnections[receiverIDStr]
					connMutex.Unlock()

					if online {
						receiverConn.WriteJSON(wsMsg)
					}

					// Echo back to sender for confirmation
					conn.WriteJSON(wsMsg)
				}
			}
		})

		return handler(c)
	})
}
