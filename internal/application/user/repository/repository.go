package repository

import (
	"context"

	"cats-social/internal/domain"
)

type AuthRepositoryContract interface {
	Register(ctx context.Context, user domain.User) error
}
