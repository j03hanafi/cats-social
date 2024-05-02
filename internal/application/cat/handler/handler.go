package handler

import (
	"github.com/gofiber/fiber/v2"

	"cats-social/internal/domain"
)

type catHandler struct{}

func NewCatHandler(router fiber.Router, jwtMiddleware fiber.Handler) {
	handler := catHandler{}

	catRouter := router.Group("/cat")

	catRouter.Use(jwtMiddleware)
	catRouter.Get("", handler.Get)
}

func (h catHandler) Get(c *fiber.Ctx) error {
	user := c.Locals(domain.UserFromToken).(domain.User)
	return c.JSON(fiber.Map{
		"message": user,
	})
}
