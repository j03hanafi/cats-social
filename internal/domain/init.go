package domain

import (
	"github.com/gofiber/fiber/v2"

	"cats-social/common/configs"
	"cats-social/internal/domain/info"
)

func New(server *fiber.App) {
	v1 := server.Group(configs.Runtime.API.BaseURL)

	info.NewModule(v1)
}
