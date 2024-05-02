package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"cats-social/common/configs"
	"cats-social/internal/domain"
)

type AccessTokenClaims struct {
	User domain.User `json:"user"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(u domain.User) (string, error) {
	callerInfo := "[security.GenerateAccessToken]"
	l := zap.L().With(zap.String("caller", callerInfo))

	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(configs.Runtime.API.JWT.Expire) * time.Second)

	claims := AccessTokenClaims{
		User: u,
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
