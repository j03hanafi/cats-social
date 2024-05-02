package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"cats-social/common/logger"
	"cats-social/internal/domain"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (a AuthRepository) Register(ctx context.Context, dUser domain.User) error {
	callerInfo := "[AuthRepository.Register]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	mUser := user{
		ID:        dUser.ID,
		Email:     dUser.Email,
		Name:      dUser.Name,
		Password:  dUser.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: sql.NullTime{
			Valid: false,
		},
	}

	tx, err := a.db.Begin(ctx)
	if err != nil {
		l.Error("failed to begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	insertQuery := `INSERT INTO users (id, email, name, password, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.Exec(
		ctx,
		insertQuery,
		mUser.ID,
		mUser.Email,
		mUser.Name,
		mUser.Password,
		mUser.CreatedAt,
		mUser.UpdatedAt,
		mUser.DeletedAt,
	)
	if err != nil {
		l.Error("failed to insert user", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		l.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

var _ AuthRepositoryContract = (*AuthRepository)(nil)
