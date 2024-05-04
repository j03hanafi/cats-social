package internal

import (
	"fmt"

	"go.uber.org/zap"

	"cats-social/common/configs"
	"cats-social/common/logger"
	"cats-social/internal/server"
)

func Run() {
	callerInfo := "[internal.Run]"

	// Load configs
	err := configs.NewConfig()
	if err != nil {
		fmt.Printf("%s failed to load config: %v\n", callerInfo, err)
		return
	}

	// Initialize logger
	l := logger.Get()
	defer func() {
		_ = l.Sync()
	}()
	zap.ReplaceGlobals(l)

	l.Debug("config loaded", zap.Any("config", configs.Runtime))

	server.Run()
}
