package info

import (
	"github.com/gofiber/fiber/v2"

	"cats-social/internal/application/info/handler"
)

func NewModule(router fiber.Router) {
	handler.NewInfoHandler(router)
}
