package cat

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"cats-social/common/configs"
	"cats-social/internal/application/cat/handler"
	catRepo "cats-social/internal/application/cat/repository"
	"cats-social/internal/application/cat/service"
	matchRepo "cats-social/internal/application/match/repository"
)

func NewModule(router fiber.Router, db *pgxpool.Pool, jwtMiddleware fiber.Handler) {
	ctxTimeout := time.Duration(configs.Runtime.App.ContextTimeout) * time.Second

	catRepository := catRepo.NewCatRepository(db)
	matchRepository := matchRepo.NewMatchRepository(db)
	catService := service.NewCatService(ctxTimeout, catRepository, matchRepository)
	handler.NewCatHandler(router, jwtMiddleware, catService)
}
