package handlers

import (
	"net/http"
	"putra4648/my-chat-app/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
	uuid "github.com/gofrs/uuid/v5"
)

type dashboardHandler struct {
	userService *services.UserService
	chatService *services.ChatService
}

func NewDashboardHandler(userService *services.UserService, chatService *services.ChatService) RouteHandler {
	return &dashboardHandler{
		userService: userService,
		chatService: chatService,
	}
}

func (h *dashboardHandler) Register(app *fiber.App) {
	app.Get("/proxy-avatar", func(c fiber.Ctx) error {
		url := "https://i.pravatar.cc/150"
		email := c.Query("email")
		if email != "" {
			url += "?u=" + email
		}

		// Ambil gambar dari eksternal
		resp, err := http.Get(url)
		if err != nil {
			return c.Status(500).SendString("Gagal ambil gambar")
		}
		defer resp.Body.Close()

		// Copy body gambar ke response Fiber
		c.Set("Content-Type", "image/jpeg")
		c.Set("Access-Control-Allow-Origin", "*") // Tambahkan CORS manual
		return c.SendStream(resp.Body)
	})

	app.Get("/", func(c fiber.Ctx) error {
		return c.Redirect().To("/dashboard")
	})

	app.Get("/dashboard", func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		csrf := csrf.TokenFromContext(c)
		userIDStr := sess.Get("user_id").(string)
		userID, _ := uuid.FromString(userIDStr)

		users, err := h.userService.GetUsersWithoutUserLogin(c.Context(), userIDStr)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		conversations, err := h.chatService.GetUserConversations(c.Context(), userID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Render("dashboard", fiber.Map{
			"Title":         "Dashboard",
			"BodyClass":     "chat-page",
			"Users":         users,
			"Conversations": conversations,
			"CurrentUserID": userIDStr,
			"_csrf":         csrf,
		})
	})

	app.Get("/conversations/:id/messages", func(c fiber.Ctx) error {
		convID, err := uuid.FromString(c.Params("id"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid conversation ID")
		}

		messages, err := h.chatService.GetConversationMessages(c.Context(), convID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(messages)
	})

	app.Post("/conversations/private/:receiverID", func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		senderID, _ := uuid.FromString(sess.Get("user_id").(string))
		receiverID, err := uuid.FromString(c.Params("receiverID"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid receiver ID")
		}

		convID, err := h.chatService.GetOrCreatePrivateConversation(c.Context(), senderID, receiverID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{"conversation_id": convID})
	})
}
