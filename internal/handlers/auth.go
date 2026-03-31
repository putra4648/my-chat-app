package handlers

import (
	"putra4648/my-chat-app/internal/models"
	"putra4648/my-chat-app/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
	service *services.UserService
}

func NewAuthHandler(service *services.UserService) RouteHandler {
	return &authHandler{service: service}
}

func (h *authHandler) Register(app *fiber.App) {
	app.Get("/login", func(c fiber.Ctx) error {
		token := csrf.TokenFromContext(c)
		return c.Render("login", fiber.Map{
			"_csrf": token,
			"BodyClass": "auth-page",
		})
	})

	app.Post("/login", func(c fiber.Ctx) error {
		sess := session.FromContext(c)

		email := c.FormValue("email")
		password := c.FormValue("password")

		user, err := h.service.GetUserByEmail(c.Context(), email)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
		}

		if err = sess.Regenerate(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Session error")
		}

		sess.Set("user_id", user.ID.String())
		sess.Set("user_email", user.Email)
		sess.Set("user_username", user.Username)
		sess.Set("authenticated", true)

		return c.Redirect().To("/dashboard")
	})

	app.Get("/register", func(c fiber.Ctx) error {
		token := csrf.TokenFromContext(c)
		return c.Render("register", fiber.Map{
			"_csrf": token,
			"BodyClass": "auth-page",
		})
	})

	app.Post("/register", func(c fiber.Ctx) error {

		username := c.FormValue("username")
		email := c.FormValue("email")
		password := c.FormValue("password")

		var user models.User
		user.Username = username
		user.Email = email
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate password hash")
		}
		user.PasswordHash = string(hash)

		_, err = h.service.CreateUser(c.Context(), &user)
		if err != nil {
			log.Error(err)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Redirect().To("/login")
	})

	app.Post("/logout", func(c fiber.Ctx) error {
		sess := session.FromContext(c)

		// Complete session reset (clears all data + new session ID)
		if err := sess.Reset(); err != nil {
			return c.Status(500).SendString("Session error")
		}

		return c.Redirect().To("/")
	})
}
