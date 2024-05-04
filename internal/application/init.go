package application

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"cats-social/common/configs"
	"cats-social/internal/application/cat"
	"cats-social/internal/application/info"
	"cats-social/internal/application/match"
	"cats-social/internal/application/user"
)

func New(server *fiber.App, db *pgxpool.Pool, jwtMiddleware fiber.Handler) {
	v1 := server.Group(configs.Runtime.API.BaseURL)

	info.NewModule(v1, db)
	user.NewModule(v1, db)
	cat.NewModule(v1, db, jwtMiddleware)
	match.NewModule(v1, db, jwtMiddleware)
}
