package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bytedance/sonic"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"cats-social/common/configs"
	"cats-social/common/database"
	"cats-social/internal/application"
)

const (
	localEnv   = "local"
	sonicParse = "sonic"
)

func Run() {
	callerInfo := "[server.Run]"
	log := zap.L().With(zap.String("caller", callerInfo))

	db, err := database.NewPGConn()
	defer db.Close()
	if err != nil {
		log.Panic("Failed to connect to database", zap.Error(err))
	}

	serverTimeout := time.Duration(configs.Runtime.API.Timeout) * time.Second
	serverConfig := fiber.Config{
		AppName:            configs.Runtime.App.Name,
		DisableDefaultDate: true,
		DisableKeepalive:   true,
		EnablePrintRoutes:  true,
		JSONDecoder:        json.Unmarshal,
		JSONEncoder:        json.Marshal,
		ReadTimeout:        serverTimeout,
		ReduceMemoryUsage:  true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := http.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return ctx.Status(code).JSON(e)
		},
	}

	if configs.Runtime.App.Env != localEnv {
		serverConfig.Prefork = true
	}

	if configs.Runtime.API.Parser == sonicParse {
		serverConfig.JSONDecoder = sonic.Unmarshal
		serverConfig.JSONEncoder = sonic.Marshal
	}

	app := fiber.New(serverConfig)
	setMiddlewares(app)
	application.New(app, db, jwtMiddleware())
	log.Debug("Server Config", zap.Any("Config", app.Config()))

	go func() {
		addr := fmt.Sprintf("%s:%d", configs.Runtime.App.Host, configs.Runtime.App.Port)
		if err := app.Listen(addr); err != nil {
			log.Panic("Server Error", zap.Error(err))
		}
	}()

	log.Info("Server is starting...")

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Info("shutting down gracefully, press Ctrl+C again to force")

	err = app.ShutdownWithTimeout(serverTimeout)
	if err != nil {
		log.Panic("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server was successful shutdown")
}
