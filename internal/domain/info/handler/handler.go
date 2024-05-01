package handler

import (
	"github.com/gofiber/fiber/v2"

	"cats-social/common/configs"
)

type infoHandler struct{}

func NewInfoHandler(server fiber.Router) {
	handler := infoHandler{}

	route := server.Group("/info")

	route.Get("/version", handler.Version)
}

func (h infoHandler) Version(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"version": configs.Runtime.App.Version,
	})
}
