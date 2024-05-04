package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"cats-social/common/id"
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

func (a AuthRepository) Create(ctx context.Context, dUser domain.User) (domain.User, error) {
	callerInfo := "[AuthRepository.Create]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	mUser := user{
		ID:        id.New(),
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
		return dUser, err
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			l.Error("user already exists", zap.Error(err))
			return dUser, domain.DuplicateEmailError
		}
		l.Error("failed to insert user", zap.Error(err))
		return dUser, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		l.Error("failed to commit transaction", zap.Error(err))
		return dUser, err
	}

	dUser.ID = mUser.ID
	return dUser, nil
}

func (a AuthRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	callerInfo := "[AuthRepository.GetByEmail]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	var mUser user
	query := `SELECT id, email, name, password FROM users WHERE email = $1`
	err := a.db.QueryRow(ctx, query, email).Scan(&mUser.ID, &mUser.Email, &mUser.Name, &mUser.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			l.Error("user not found", zap.Error(err))
			return domain.User{}, domain.UserNotFoundError
		}
		l.Error("failed to get user", zap.Error(err))
		return domain.User{}, err
	}

	dUser := domain.User{
		ID:       mUser.ID,
		Email:    mUser.Email,
		Name:     mUser.Name,
		Password: mUser.Password,
	}

	return dUser, nil
}

func (a AuthRepository) Get(ctx context.Context, userID ulid.ULID) (domain.User, error) {
	callerInfo := "[AuthRepository.Get]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	var mUser user
	query := `SELECT email, name, created_at FROM users WHERE id = $1`
	err := a.db.QueryRow(ctx, query, userID).Scan(&mUser.Email, &mUser.Name, &mUser.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			l.Error("user not found", zap.Error(err))
			return domain.User{}, domain.UserNotFoundError
		}
		l.Error("failed to get user", zap.Error(err))
		return domain.User{}, err
	}

	dUser := domain.User{
		ID:        userID,
		Name:      mUser.Name,
		Email:     mUser.Email,
		CreatedAt: mUser.CreatedAt,
	}

	return dUser, nil
}

var _ AuthRepositoryContract = (*AuthRepository)(nil)
