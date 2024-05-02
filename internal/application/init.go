package application

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"cats-social/common/configs"
	"cats-social/internal/application/info"
	"cats-social/internal/application/user"
)

func New(server *fiber.App, db *pgxpool.Pool) {
	v1 := server.Group(configs.Runtime.API.BaseURL)

	info.NewModule(v1, db)
	user.NewModule(v1, db)
}
