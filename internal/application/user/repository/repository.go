package repository

import (
	"context"

	"github.com/oklog/ulid/v2"

	"cats-social/internal/domain"
)

type AuthRepositoryContract interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Get(ctx context.Context, userID ulid.ULID) (domain.User, error)
}
