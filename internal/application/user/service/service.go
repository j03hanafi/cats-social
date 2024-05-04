package service

import (
	"context"

	"cats-social/internal/domain"
)

type AuthServiceContract interface {
	Register(ctx context.Context, user domain.User) (domain.User, error)
	GenerateToken(ctx context.Context, user domain.User) (string, error)
	Login(ctx context.Context, user domain.User) (domain.User, error)
}
