package server

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"

	"cats-social/common/id"
	"cats-social/common/logger"
)

const (
	requestId = "requestId"
)

func setMiddlewares(app *fiber.App) {
	app.Use(compressionMiddleware())
	app.Use(recoveryMiddleware())
	app.Use(zapMiddleware())
	app.Use(requestIDMiddleware())
	app.Use(loggerMiddleware())
	app.Use(cacheMiddleware())
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
