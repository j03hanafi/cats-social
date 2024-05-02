package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"cats-social/common/configs"
)

func setConsoleLogger() (zapcore.Core, []zap.Option) {
	writer := zapcore.AddSync(os.Stdout)

	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	encoder := zapcore.NewConsoleEncoder(config)

	logLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	if !configs.Runtime.API.DebugMode {
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	options := append([]zap.Option{}, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))

	return zapcore.NewCore(encoder, writer, logLevel), options
}
