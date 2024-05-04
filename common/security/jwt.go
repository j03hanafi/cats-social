package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"

	"cats-social/common/configs"
	"cats-social/internal/domain"
)

type AccessTokenClaims struct {
	User user `json:"user"`
	jwt.RegisteredClaims
}

type user struct {
	ID    ulid.ULID `json:"id"`
	Email string    `json:"email"`
	Name  string    `json:"name"`
}

func GenerateAccessToken(u domain.User) (string, error) {
	callerInfo := "[security.GenerateAccessToken]"
	l := zap.L().With(zap.String("caller", callerInfo))

	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(configs.Runtime.API.JWT.Expire) * time.Second)

	claims := AccessTokenClaims{
		User: user{
			ID:    u.ID,
			Email: u.Email,
			Name:  u.Name,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ExpiresAt: jwt.NewNumericDate(tokenExp),
			NotBefore: jwt.NewNumericDate(currentTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(configs.Runtime.API.JWT.JWTSecret))
	if err != nil {
		l.Error("Error signing token",
			zap.Error(err),
		)
		return "", err
	}

	return signedString, nil
}
