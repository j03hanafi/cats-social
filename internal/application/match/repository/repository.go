package repository

import (
	"context"

	"github.com/oklog/ulid/v2"

	"cats-social/internal/domain"
)

type MatchRepositoryContract interface {
	NewMatch(ctx context.Context, match domain.Match) (domain.Match, error)
	HasMatched(ctx context.Context, match domain.Match) (bool, error)
	GetDetailMatches(ctx context.Context, userID ulid.ULID) ([]domain.DetailMatch, error)
	Get(ctx context.Context, matchID ulid.ULID) (domain.DetailMatch, error)
	DeleteExceptApproved(ctx context.Context, userID, matchID ulid.ULID) error
	Delete(ctx context.Context, matchID ulid.ULID) error
}
