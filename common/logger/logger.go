package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxKey struct{}

var (
	once   sync.Once
	logger *zap.Logger
)

// Get initializes a zap.Logger instance if it has not been initialized
// already and returns the same instance for subsequent calls.
func Get() *zap.Logger {
	once.Do(func() {
		options := make([]zap.Option, 0)

		consoleCore, consoleOpts := setConsoleLogger()
		options = append(options, consoleOpts...)

		fileCore, fileOpts := setFileLogger()
		options = append(options, fileOpts...)

		core := zapcore.NewTee(
			consoleCore,
			fileCore,
		)

		logger = zap.New(core, options...)
	})

	return logger
}

// FromCtx returns the Logger associated with the ctx. If no logger
// is associated, the default logger is returned, unless it is nil
// in which case a disabled logger is returned.
func FromCtx(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	} else if l := logger; l != nil {
		return l
	}

	return zap.NewNop()
}

// WithCtx returns a copy of ctx with the Logger attached.
func WithCtx(ctx context.Context, l *zap.Logger) context.Context {
	if lp, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		if lp == l {
			// Do not store same logger.
			return ctx
		}
	}

	return context.WithValue(ctx, ctxKey{}, l)
}
