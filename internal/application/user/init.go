package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"cats-social/common/configs"
	"cats-social/internal/application/user/handler"
	"cats-social/internal/application/user/repository"
	"cats-social/internal/application/user/service"
)

func NewModule(router fiber.Router, db *pgxpool.Pool) {
	ctxTimeout := time.Duration(configs.Runtime.App.ContextTimeout) * time.Second

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(ctxTimeout, authRepository)
	handler.NewAuthHandler(router, authService)
}
