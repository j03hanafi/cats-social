package internal

import (
	"fmt"

	"go.uber.org/zap"

	"cats-social/common/configs"
	"cats-social/common/logger"
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

	zap.L().Info("Configs", zap.Any("Runtime", configs.Runtime))
}
