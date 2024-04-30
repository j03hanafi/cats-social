package logger

import (
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"cats-social/common/configs"
)

const (
	logMaxSize    = 5
	logMaxAge     = 30
	logMaxBackups = 15
	logSize       = 512 * 1024

	bufferFlushInterval = time.Second

	samplerFirst      = 100
	samplerThereafter = 100
)

func setFileLogger() (zapcore.Core, []zap.Option) {
	filename := "logs/" + configs.Runtime.App.Name + ".log"
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    logMaxSize,
		MaxAge:     logMaxAge,
		MaxBackups: logMaxBackups,
		LocalTime:  true,
	})
	bufferedWriter := &zapcore.BufferedWriteSyncer{
		WS:            writer,
		Size:          logSize,
		FlushInterval: bufferFlushInterval,
	}

	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(config)

	logLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	options := make([]zap.Option, 0)

	var gitRevision, goVersion string
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		for _, v := range buildInfo.Settings {
			if v.Key == "vcs.revision" {
				gitRevision = v.Value
				break
			}
		}
		goVersion = buildInfo.GoVersion
	}

	core := zapcore.NewCore(encoder, bufferedWriter, logLevel).
		With([]zap.Field{
			zap.String("gitRevision", gitRevision),
			zap.String("goVersion", goVersion),
		})

	return zapcore.NewSamplerWithOptions(core, time.Second, samplerFirst, samplerThereafter), options
}
