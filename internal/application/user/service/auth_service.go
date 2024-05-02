package service

import (
	"context"
	"time"

	"go.uber.org/zap"

	"cats-social/common/id"
	"cats-social/common/logger"
	"cats-social/common/security"
	"cats-social/internal/application/user/repository"
	"cats-social/internal/domain"
)

type AuthService struct {
	authRepository repository.AuthRepositoryContract
	contextTimeout time.Duration
}

func NewAuthService(timeout time.Duration, authRepository repository.AuthRepositoryContract) *AuthService {
	authService := &AuthService{
		authRepository: authRepository,
		contextTimeout: timeout,
	}

	return authService
}

func (a AuthService) Register(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	callerInfo := "[AuthService.Register]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	password, err := security.HashPassword(user.Password)
	if err != nil {
		l.Error("error hashing password",
			zap.Error(err),
		)
		return domain.User{}, err
	}

	user.Password = password
	user.ID = id.New()

	if err = a.authRepository.Register(ctx, user); err != nil {
		l.Error("error register user",
			zap.Error(err),
		)
		return domain.User{}, err
	}

	return user, nil
}

var _ AuthServiceContract = (*AuthService)(nil)
