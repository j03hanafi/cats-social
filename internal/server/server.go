package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"cats-social/common/configs"
	"cats-social/internal/domain"
)

const (
	localEnv = "local"
)

func Run() {
	callerInfo := "[server.Run]"
	log := zap.L().With(zap.String("caller", callerInfo))

	serverTimeout := time.Duration(configs.Runtime.API.Timeout) * time.Second
	serverConfig := fiber.Config{
		AppName:            configs.Runtime.App.Name,
		DisableDefaultDate: true,
		DisableKeepalive:   true,
		EnablePrintRoutes:  true,
		JSONDecoder:        sonic.Unmarshal,
		JSONEncoder:        sonic.Marshal,
		ReadTimeout:        serverTimeout,
	}

	if configs.Runtime.App.Env != localEnv {
		serverConfig.Prefork = true
	}

	app := fiber.New(serverConfig)
	setMiddlewares(app)
	domain.New(app)

	go func() {
		addr := fmt.Sprintf("%s:%d", configs.Runtime.App.Host, configs.Runtime.App.Port)
		if err := app.Listen(addr); err != nil {
			log.DPanic("Server Error", zap.Error(err))
		}
	}()

	log.Info("Server is starting...")

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Info("shutting down gracefully, press Ctrl+C again to force")

	err := app.ShutdownWithTimeout(serverTimeout)
	if err != nil {
		log.DPanic("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server was successful shutdown")
}
