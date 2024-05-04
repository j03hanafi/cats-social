package match

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"cats-social/common/configs"
	catRepo "cats-social/internal/application/cat/repository"
	"cats-social/internal/application/match/handler"
	matchRepo "cats-social/internal/application/match/repository"
	"cats-social/internal/application/match/service"
	userRepo "cats-social/internal/application/user/repository"
)

func NewModule(router fiber.Router, db *pgxpool.Pool, jwtMiddleware fiber.Handler) {
	ctxTimeout := time.Duration(configs.Runtime.App.ContextTimeout) * time.Second

	catRepository := catRepo.NewCatRepository(db)
	userRepository := userRepo.NewAuthRepository(db)
	matchRepository := matchRepo.NewMatchRepository(db)
	matchService := service.NewMatchService(ctxTimeout, matchRepository, catRepository, userRepository)
	handler.NewMatchHandler(router, jwtMiddleware, matchService)
}
