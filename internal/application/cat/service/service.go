package service

import (
	"context"

	"github.com/oklog/ulid/v2"

	"cats-social/internal/domain"
)

type CatServiceContract interface {
	AddCat(ctx context.Context, cat domain.Cat) (domain.Cat, error)
	ListCats(ctx context.Context, userID ulid.ULID, query domain.QueryParam) ([]domain.Cat, error)
	UpdateCat(ctx context.Context, cat domain.Cat) (domain.Cat, error)
	DeleteCat(ctx context.Context, cat domain.Cat) error
}
