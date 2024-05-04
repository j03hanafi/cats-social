package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"cats-social/common/configs"
)

type infoHandler struct {
	db *pgxpool.Pool
}

func NewInfoHandler(router fiber.Router, db *pgxpool.Pool) {
	handler := infoHandler{
		db: db,
	}

	infoRouter := router.Group("/info")

	infoRouter.Get("/version", handler.Version)
	infoRouter.Get("/health", handler.Health)
}

func (h infoHandler) Version(c *fiber.Ctx) error {
	versionInfo := version{
		Version: configs.Runtime.App.Version,
	}

	res := baseResponse{
		Message: "API version",
		Data:    versionInfo,
	}

	return c.JSON(res)
}

func (h infoHandler) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(configs.Runtime.App.ContextTimeout)*time.Second,
	)
	defer cancel()

	var errs []error
	conns := h.db.AcquireAllIdle(ctx)
	for _, conn := range conns {
		if err := conn.Ping(ctx); err != nil {
			errs = append(errs, err)
		}
		conn.Release()
	}

	if len(errs) > 0 {
		res := baseResponse{
			Message: "API is unhealthy",
			Data: fiber.Map{
				"errors": errs,
			},
		}
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	res := baseResponse{
		Message: "API is up and running",
		Data:    fiber.Map{},
	}

	dbData := fiber.Map{
		"status":    "connected",
		"totalIdle": len(conns),
		"stat":      h.db.Stat().MaxConns(),
	}

	if configs.Runtime.API.DebugMode {
		res.Data = dbData
	}

	return c.JSON(res)
}
