package info

import (
	"github.com/gofiber/fiber/v2"

	"cats-social/internal/domain/info/handler"
)

func NewModule(server fiber.Router) {
	handler.NewInfoHandler(server)
}
