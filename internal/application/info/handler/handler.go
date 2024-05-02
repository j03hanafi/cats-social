package handler

import (
	"github.com/gofiber/fiber/v2"

	"cats-social/common/configs"
)

type infoHandler struct{}

func NewInfoHandler(router fiber.Router) {
	handler := infoHandler{}

	infoRouter := router.Group("/info")

	infoRouter.Get("/version", handler.Version)
}

func (h infoHandler) Version(ctx *fiber.Ctx) error {
	versionInfo := version{
		Version: configs.Runtime.App.Version,
	}

	res := baseResponse{
		Message: "API version",
		Data:    versionInfo,
	}

	return ctx.JSON(res)
}
