package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	pgxZap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"go.uber.org/zap"

	"cats-social/common/configs"
)

func NewPGConn() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(configs.Runtime.App.ContextTimeout)*time.Second,
	)
	defer cancel()

	callerInfo := "[database.NewPGConn]"
	l := zap.L().With(zap.String("caller", callerInfo))

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
		configs.Runtime.DB.Host,
		configs.Runtime.DB.Port,
		configs.Runtime.DB.Username,
		configs.Runtime.DB.Password,
		configs.Runtime.DB.Name,
		strings.Join(configs.Runtime.DB.Params, " "),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		l.Error("error parsing database config",
			zap.Error(err),
		)
		return nil
	}

	poolLog := pgxZap.NewLogger(zap.L())
	poolTracer := &tracelog.TraceLog{
		Logger:   poolLog,
		LogLevel: tracelog.LogLevelDebug,
	}
	if !configs.Runtime.API.DebugMode {
		poolTracer.LogLevel = tracelog.LogLevelNone
	}
	config.ConnConfig.Tracer = poolTracer

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		l.Error("error creating database pool",
			zap.Error(err),
		)
		return nil
	}

	if err = pool.Ping(ctx); err != nil {
		l.Error("error pinging database",
			zap.Error(err),
		)
		return nil
	} else {
		l.Info("connected to database")
	}

	return pool
}
