package cat

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"cats-social/internal/application/cat/handler"
)

func NewModule(router fiber.Router, db *pgxpool.Pool, jwtMiddleware fiber.Handler) {
	// ctxTimeout := time.Duration(configs.Runtime.App.ContextTimeout) * time.Second

	handler.NewCatHandler(router, jwtMiddleware)
}
