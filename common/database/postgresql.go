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

func NewPGConn() (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(configs.Runtime.App.ContextTimeout)*time.Second,
	)
	defer cancel()

	callerInfo := "[database.NewPGConn]"
	l := zap.L().With(zap.String("caller", callerInfo))

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s&pool_max_conns=%d",
		configs.Runtime.DB.Username,
		configs.Runtime.DB.Password,
		configs.Runtime.DB.Host,
		configs.Runtime.DB.Port,
		configs.Runtime.DB.Name,
		strings.Join(configs.Runtime.DB.Params, " "),
		configs.Runtime.DB.MaxConnPool,
	)

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		l.Error("error parsing database config",
			zap.Error(err),
		)
		return nil, err
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
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		l.Error("error pinging database",
			zap.Error(err),
		)
		return nil, err
	} else {
		l.Info("connected to database")
	}

	l.Debug("Database Config", zap.Any("Config", pool.Config().ConnString()))
	return pool, nil
}
