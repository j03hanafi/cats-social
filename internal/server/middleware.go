package server

import (
	"net/http"

	"github.com/gofiber/contrib/fiberzap/v2"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"cats-social/common/configs"
	"cats-social/common/id"
	"cats-social/common/logger"
	"cats-social/common/security"
	"cats-social/internal/domain"
)

const (
	requestId   = "requestId"
	accessToken = "accessToken"
)

func setMiddlewares(app *fiber.App) {
	app.Use(compressionMiddleware())
	app.Use(recoveryMiddleware())
	app.Use(zapMiddleware())
	app.Use(requestIDMiddleware())
	app.Use(loggerMiddleware())
	// app.Use(cacheMiddleware())
	app.Use(eTagMiddleware())
}

func compressionMiddleware() fiber.Handler {
	return compress.New()
}

func recoveryMiddleware() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
	})
}

func zapMiddleware() fiber.Handler {
	return fiberzap.New(fiberzap.Config{
		Logger: zap.L(),
		Fields: []string{
			"latency",
			"time",
			"requestId",
			"pid",
			"status",
			"method",
			"path",
			"queryParams",
			"ip",
			"ua",
			"resBody",
			"error",
		},
	})
}

func requestIDMiddleware() fiber.Handler {
	return requestid.New(requestid.Config{
		Generator: func() string {
			return id.New().String()
		},
		ContextKey: requestId,
	})
}

func loggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		l := zap.L().With(zap.String(requestId, c.Locals(requestId).(string)))
		ctx = logger.WithCtx(ctx, l)
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func cacheMiddleware() fiber.Handler {
	return cache.New(cache.Config{
		CacheControl: true,
	})
}

func eTagMiddleware() fiber.Handler {
	return etag.New()
}

func jwtMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(configs.Runtime.API.JWT.JWTSecret),
		},
		Claims:     &security.AccessTokenClaims{},
		ContextKey: accessToken,
		SuccessHandler: func(c *fiber.Ctx) error {
			claims := c.Locals(accessToken).(*jwt.Token).Claims.(*security.AccessTokenClaims)

			user := domain.User{
				ID:    claims.User.ID,
				Email: claims.User.Email,
				Name:  claims.User.Name,
			}
			c.Locals(domain.UserFromToken, user)
			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		},
	})
}
