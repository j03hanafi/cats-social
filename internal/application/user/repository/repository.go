package repository

import (
	"context"

	"cats-social/internal/domain"
)

type AuthRepositoryContract interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}
