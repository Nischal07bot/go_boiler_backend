package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/Nischal07bot/go_boiler_backend/internal/config"
	loggerConfig "github.com/Nischal07bot/go_boiler_backend/internal/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type Database struct {
	Pool *pgxpool.Pool //this will be exported to be used in other packages do convention capital letters
	log  *zerolog.Logger
}

// allows chaining multiple tracers together to create a single tracer that can be used for logging and monitoring purposes. This is useful for applications that need to support multiple tracing systems or want to combine different types of tracers for more comprehensive observability.
type multitracer struct {
	tracers []any
}

func (mt *multitracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryStart(context.Context, *pgx.Conn, pgx.TraceQueryStartData) context.Context
		}); ok {
			ctx = t.TraceQueryStart(ctx, conn, data)
		}
	}

	return ctx
}

func (mt *multitracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	for _, tracer := range mt.tracers {
		if t, ok := tracer.(interface {
			TraceQueryEnd(context.Context, *pgx.Conn, pgx.TraceQueryEndData)
		}); ok {
			t.TraceQueryEnd(ctx, conn, data)
		}
	}
}

type pgxTraceLogAdapter struct {
	log zerolog.Logger
}

func (a *pgxTraceLogAdapter) Log(_ context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	event := a.log.WithLevel(zerolog.InfoLevel)

	switch int(level) {
	case 6:
		event = a.log.WithLevel(zerolog.DebugLevel)
	case 5:
		event = a.log.WithLevel(zerolog.DebugLevel)
	case 4:
		event = a.log.WithLevel(zerolog.InfoLevel)
	case 3:
		event = a.log.WithLevel(zerolog.WarnLevel)
	case 2:
		event = a.log.WithLevel(zerolog.ErrorLevel)
	default:
		event = a.log.WithLevel(zerolog.InfoLevel)
	}

	for key, value := range data {
		event = event.Interface(key, value)
	}

	event.Msg(msg)
}

const DatabasePingTimeout = 10

func New(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerConfig.LoggerService) (*Database, error) {
	hostPort := net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port))

	// URL-encode the password
	encodedPassword := url.QueryEscape(cfg.Database.Password)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.Database.User,
		encodedPassword,
		hostPort,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	_ = loggerService

	if cfg.Primary.Env == "local" {
		globalLevel := logger.GetLevel()
		pgxLogger := loggerConfig.NewPgxLogger(globalLevel)
		// Chain tracers - New Relic first, then local logging
		if pgxPoolConfig.ConnConfig.Tracer != nil {
			// If New Relic tracer exists, create a multi-tracer
			localTracer := &tracelog.TraceLog{
				Logger:   &pgxTraceLogAdapter{log: pgxLogger},
				LogLevel: tracelog.LogLevel(loggerConfig.Getpgxtraceloglevel(globalLevel)),
			}
			pgxPoolConfig.ConnConfig.Tracer = &multitracer{
				tracers: []any{pgxPoolConfig.ConnConfig.Tracer, localTracer},
			}
		} else {
			pgxPoolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
				Logger:   &pgxTraceLogAdapter{log: pgxLogger},
				LogLevel: tracelog.LogLevel(loggerConfig.Getpgxtraceloglevel(globalLevel)),
			}
		}
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	database := &Database{
		Pool: pool,
		log:  logger,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DatabasePingTimeout*time.Second)
	defer cancel()
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info().Msg("connected to the database")

	return database, nil
}

func (db *Database) Close() error {
	db.log.Info().Msg("database connection closed")
	db.Pool.Close()
	return nil
}