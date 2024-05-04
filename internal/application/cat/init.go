package cat

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"cats-social/common/configs"
	"cats-social/internal/application/cat/handler"
	"cats-social/internal/application/cat/repository"
	"cats-social/internal/application/cat/service"
)

func NewModule(router fiber.Router, db *pgxpool.Pool, jwtMiddleware fiber.Handler) {
	ctxTimeout := time.Duration(configs.Runtime.App.ContextTimeout) * time.Second

	catRepository := repository.NewCatRepository(db)
	catService := service.NewCatService(ctxTimeout, catRepository)
	handler.NewCatHandler(router, jwtMiddleware, catService)
}
