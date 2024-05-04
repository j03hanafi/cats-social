package service

import (
	"context"

	"github.com/oklog/ulid/v2"

	"cats-social/internal/domain"
)

type MatchServiceContract interface {
	NewMatch(ctx context.Context, match domain.Match, userID ulid.ULID) (domain.Match, error)
	GetMatch(ctx context.Context, userID ulid.ULID) ([]domain.DetailMatch, error)
	ApproveMatch(ctx context.Context, matchID, userID ulid.ULID) error
	RejectMatch(ctx context.Context, matchID, userID ulid.ULID) error
	DeleteMatch(ctx context.Context, matchID, userID ulid.ULID) error
}
