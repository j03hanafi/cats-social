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

	callerInfo := "[AuthService.Create]"
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

	if err = a.authRepository.Create(ctx, user); err != nil {
		l.Error("error register user",
			zap.Error(err),
		)
		return domain.User{}, err
	}

	return user, nil
}

func (a AuthService) GenerateToken(ctx context.Context, user domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	callerInfo := "[AuthService.GenerateToken]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	user.Password = ""
	token, err := security.GenerateAccessToken(user)
	if err != nil {
		l.Error("error generating token",
			zap.Error(err),
		)
		return "", err
	}

	return token, nil
}

func (a AuthService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	callerInfo := "[AuthService.Login]"
	l := logger.FromCtx(ctx).With(zap.String("caller", callerInfo))

	userData, err := a.authRepository.GetByEmail(ctx, user.Email)
	if err != nil {
		l.Error("error get user by email",
			zap.Error(err),
		)
		return domain.User{}, err
	}

	if err = security.ComparePasswords(userData.Password, user.Password); err != nil {
		l.Error("error compare password",
			zap.Error(err),
		)
		return domain.User{}, domain.InvalidPassword
	}

	return userData, nil
}

var _ AuthServiceContract = (*AuthService)(nil)
