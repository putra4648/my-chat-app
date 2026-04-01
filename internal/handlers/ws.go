package handlers

import (
	"context"
	"encoding/json"
	"putra4648/my-chat-app/internal/models"
	"putra4648/my-chat-app/internal/services"
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	uuid "github.com/gofrs/uuid/v5"
)

type WsMessage struct {
	Type           string    `json:"type"` // "text", "join", etc.
	ConversationID uuid.UUID `json:"conversation_id"`
	ReceiverID     uuid.UUID `json:"receiver_id"`
	Content        string    `json:"content"`
	SenderID       uuid.UUID `json:"sender_id"`
}

type Hub struct {
	chatService *services.ChatService
	clients     map[string]*Client // userID -> Client
	broadcast   chan WsMessage
	register    chan *Client
	unregister  chan *Client
	mutex       sync.Mutex
}

func NewHub(chatService *services.ChatService) *Hub {
	return &Hub{
		chatService: chatService,
		clients:     make(map[string]*Client),
		broadcast:   make(chan WsMessage),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.userID] = client
			h.mutex.Unlock()
		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mutex.Unlock()
		case message := <-h.broadcast:
			go func(msg WsMessage) {
				dbMsg := &models.Message{
					ConversationID: msg.ConversationID,
					SenderID:       &msg.SenderID,
					Content:        msg.Content,
				}
				_ = h.chatService.SaveMessage(context.Background(), dbMsg)
			}(message)

			h.mutex.Lock()
			if receiver, ok := h.clients[message.ReceiverID.String()]; ok {
				select {
				case receiver.send <- message:
				default:
					h.unregisterClient(receiver)
				}
			}
			if sender, ok := h.clients[message.SenderID.String()]; ok {
				select {
				case sender.send <- message:
				default:
					h.unregisterClient(sender)
				}
			}
			h.mutex.Unlock()
		}
	}
}

func (h *Hub) unregisterClient(client *Client) {
	delete(h.clients, client.userID)
	close(client.send)
}

type Client struct {
	hub    *Hub
	userID string
	conn   *websocket.Conn
	send   chan WsMessage
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var wsMsg WsMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}
		// Use a goroutine to avoid blocking the read loop if the hub is busy
		go func(msg WsMessage) {
			c.hub.broadcast <- msg
		}(wsMsg)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		message, ok := <-c.send
		if !ok {
			return
		}
		if err := c.conn.WriteJSON(message); err != nil {
			return
		}
	}
}

type wsHandler struct {
	hub *Hub
}

func NewWSHandler(chatService *services.ChatService) RouteHandler {
	hub := NewHub(chatService)
	go hub.Run()
	return &wsHandler{hub: hub}
}

func (h *wsHandler) Register(app *fiber.App) {
	app.Get("/ws", func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		userID := sess.Get("user_id")
		if userID == nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
		}
		c.Locals("user_id", userID)
		return c.Next()
	}, websocket.New(func(c *websocket.Conn) {
		userID := c.Locals("user_id")
		if userID == nil {
			_ = c.Close()
			return
		}
		userIDStr := userID.(string)

		client := &Client{
			hub:    h.hub,
			userID: userIDStr,
			conn:   c,
			send:   make(chan WsMessage, 256),
		}
		h.hub.register <- client

		go client.writePump()
		client.readPump()
	}))
}
