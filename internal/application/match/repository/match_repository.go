package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"cats-social/common/id"
	"cats-social/common/logger"
	"cats-social/internal/domain"
)

type MatchRepository struct {
	db *pgxpool.Pool
}

func NewMatchRepository(db *pgxpool.Pool) *MatchRepository {
	return &MatchRepository{
		db: db,
	}
}

func (m MatchRepository) NewMatch(ctx context.Context, dMatch domain.Match) (domain.Match, error) {
	callerInfo := "[MatchRepository.NewMatch]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	tx, err := m.db.Begin(ctx)
	if err != nil {
		l.Error("error starting transaction",
			zap.Error(err),
		)
		return dMatch, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	mMatch := match{
		ID:         id.New(),
		MatchCatID: dMatch.MatchCatID,
		UserCatID:  dMatch.UserCatID,
		Message:    dMatch.Message,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DeletedAt: sql.NullTime{
			Valid: false,
		},
	}

	insertQuery := `INSERT INTO matches (id, match_cat_id, user_cat_id, message, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = tx.Exec(
		ctx,
		insertQuery,
		mMatch.ID,
		mMatch.MatchCatID,
		mMatch.UserCatID,
		mMatch.Message,
		mMatch.CreatedAt,
		mMatch.UpdatedAt,
		mMatch.DeletedAt,
	)
	if err != nil {
		l.Error("error inserting data",
			zap.Error(err),
		)
		return dMatch, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		l.Error("error committing transaction",
			zap.Error(err),
		)
		return dMatch, err
	}

	dMatch.ID = mMatch.ID
	return dMatch, nil
}

func (m MatchRepository) HasMatched(ctx context.Context, dMatch domain.Match) (bool, error) {
	callerInfo := "[MatchRepository.FindMatch]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	var mMatch match
	query := `SELECT id, match_cat_id, user_cat_id, message, created_at, updated_at, deleted_at FROM matches WHERE ((match_cat_id = $1 AND user_cat_id = $2) OR (match_cat_id = $3 AND user_cat_id = $4)) AND deleted_at IS NULL`
	err := m.db.QueryRow(ctx, query, dMatch.MatchCatID, dMatch.UserCatID, dMatch.UserCatID, dMatch.MatchCatID).
		Scan(&mMatch.ID, &mMatch.MatchCatID, &mMatch.UserCatID, &mMatch.Message, &mMatch.CreatedAt, &mMatch.UpdatedAt, &mMatch.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		l.Error("error scanning data",
			zap.Error(err),
		)
		return false, err
	}

	return true, nil
}

func (m MatchRepository) GetDetailMatches(ctx context.Context, userID ulid.ULID) ([]domain.DetailMatch, error) {
	callerInfo := "[MatchRepository.GetDetailMatches]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	getMatchesQuery := `SELECT m.id, m.match_cat_id, m.user_cat_id, m.message, m.created_at, r.user_id as receiver_id, i.user_id as issuer_id
		FROM matches m 
		JOIN cats r ON m.match_cat_id = r.id
		JOIN cats i ON m.user_cat_id = i.id
		WHERE (r.user_id = $1 OR i.user_id = $2) AND m.deleted_at IS NULL ORDER BY m.created_at DESC`

	rows, err := m.db.Query(ctx, getMatchesQuery, userID, userID)
	if err != nil {
		l.Error("error getting data",
			zap.Error(err),
		)
		return nil, err
	}
	defer rows.Close()

	matches := make([]domain.DetailMatch, 0)
	for rows.Next() {
		var mMatch match
		err = rows.Scan(
			&mMatch.ID,
			&mMatch.MatchCatID,
			&mMatch.UserCatID,
			&mMatch.Message,
			&mMatch.CreatedAt,
			&mMatch.ReceiverID,
			&mMatch.IssuerID,
		)
		if err != nil {
			l.Error("error scanning data",
				zap.Error(err),
			)
			return nil, err
		}

		matches = append(matches, domain.DetailMatch{
			Match: domain.Match{
				ID:         mMatch.ID,
				MatchCatID: mMatch.MatchCatID,
				UserCatID:  mMatch.UserCatID,
				Message:    mMatch.Message,
				CreatedAt:  mMatch.CreatedAt,
			},
			Issuer: domain.User{
				ID: mMatch.IssuerID,
			},
			Receiver: domain.User{
				ID: mMatch.ReceiverID,
			},
		})
	}

	if err = rows.Err(); err != nil {
		l.Error("error iterating data",
			zap.Error(err),
		)
		return nil, err
	}

	return matches, nil
}

func (m MatchRepository) Get(ctx context.Context, matchID ulid.ULID) (domain.DetailMatch, error) {
	callerInfo := "[MatchRepository.Get]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	getMatchQuery := `SELECT m.id, m.match_cat_id, m.user_cat_id, m.message, m.created_at, r.user_id as receiver_id, i.user_id as issuer_id, m.deleted_at
		FROM matches m
		JOIN cats r ON m.match_cat_id = r.id
		JOIN cats i ON m.user_cat_id = i.id
		WHERE m.id = $1`

	var mMatch match
	err := m.db.QueryRow(ctx, getMatchQuery, matchID).Scan(
		&mMatch.ID,
		&mMatch.MatchCatID,
		&mMatch.UserCatID,
		&mMatch.Message,
		&mMatch.CreatedAt,
		&mMatch.ReceiverID,
		&mMatch.IssuerID,
		&mMatch.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.DetailMatch{}, domain.ErrMatchNotFound
		}
		l.Error("error scanning data",
			zap.Error(err),
		)
		return domain.DetailMatch{}, err
	}

	var deletedAt time.Time
	if mMatch.DeletedAt.Valid {
		deletedAt = mMatch.DeletedAt.Time
	}

	return domain.DetailMatch{
		Match: domain.Match{
			ID:         mMatch.ID,
			MatchCatID: mMatch.MatchCatID,
			UserCatID:  mMatch.UserCatID,
			Message:    mMatch.Message,
			CreatedAt:  mMatch.CreatedAt,
			DeletedAt:  deletedAt,
		},
		Issuer: domain.User{
			ID: mMatch.IssuerID,
		},
		Receiver: domain.User{
			ID: mMatch.ReceiverID,
		},
	}, nil
}

func (m MatchRepository) DeleteExceptApproved(ctx context.Context, userID, matchID ulid.ULID) error {
	callerInfo := "[MatchRepository.DeleteExceptApproved]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	tx, err := m.db.Begin(ctx)
	if err != nil {
		l.Error("error starting transaction",
			zap.Error(err),
		)
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	deleteQuery := `UPDATE matches SET deleted_at = $1
		FROM cats as r, cats as i
		WHERE matches.match_cat_id = r.id AND matches.user_cat_id = i.id
		AND (r.user_id = $2 OR i.user_id = $3) AND matches.id != $4`

	_, err = tx.Exec(ctx, deleteQuery, time.Now(), userID, userID, matchID)
	if err != nil {
		l.Error("error deleting data",
			zap.Error(err),
		)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		l.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (m MatchRepository) Delete(ctx context.Context, matchID ulid.ULID) error {
	callerInfo := "[MatchRepository.Delete]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	tx, err := m.db.Begin(ctx)
	if err != nil {
		l.Error("error starting transaction",
			zap.Error(err),
		)
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	deleteQuery := `UPDATE matches SET deleted_at = $1 WHERE id = $2`

	_, err = tx.Exec(ctx, deleteQuery, time.Now(), matchID)
	if err != nil {
		l.Error("error deleting data",
			zap.Error(err),
		)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		l.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

var _ MatchRepositoryContract = (*MatchRepository)(nil)
