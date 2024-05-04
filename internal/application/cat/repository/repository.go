package repository

import (
	"context"

	"github.com/oklog/ulid/v2"

	"cats-social/internal/domain"
)

type CatRepositoryContract interface {
	Create(ctx context.Context, cat domain.Cat) (domain.Cat, error)
	Get(ctx context.Context, userID ulid.ULID, query domain.QueryParam, withImages bool) ([]domain.Cat, error)
	Update(ctx context.Context, cat domain.Cat) (domain.Cat, error)
	Delete(ctx context.Context, catID ulid.ULID) error
}
