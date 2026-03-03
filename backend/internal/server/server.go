package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Nischal07bot/go_boiler_backend/internal/config"
	"github.com/Nischal07bot/go_boiler_backend/internal/database"
	"github.com/Nischal07bot/go_boiler_backend/internal/lib/job"
	loggerpkg "github.com/Nischal07bot/go_boiler_backend/internal/logger"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Server struct { //central server struct having all the dependencies required for the server to run
	Config        *config.Config
	Logger        *zerolog.Logger
	LoggerService *loggerpkg.LoggerService
	DB            *database.Database
	Redis         *redis.Client
	httpServer    *http.Server
	Job           *job.Jobservice
} //dependency injection create all the dependencies at the startup and store them all in the server struct  pass server to handlers that need them
/*┌─────────────────────────────────────────────────────┐
│                    Server Struct                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│   Config ──────► Environment settings               │
│   Logger ──────► Logging                            │
│   DB ──────────► PostgreSQL/MySQL                   │
│   Redis ───────► Cache + Job Queue                  │
│   Job ─────────► Background tasks (emails, etc.)    │
│   httpServer ──► HTTP API                           │
│                                                     │
└─────────────────────────────────────────────────────┘*/

func New(cfg *config.Config, logger *zerolog.Logger, loggerService *loggerpkg.LoggerService) (*Server, error) {
	db, err := database.New(cfg, logger, loggerService)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)

	}

	//redis client with New relic integration
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})

	//add new relic redis hooks if available
	if loggerService != nil && loggerService.GetApplication() != nil {
		redisClient.AddHook(nrredis.NewHook((redisClient.Options())))
	}

	//Test redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	//job service 
	jobService := job.NewJobService(logger,cfg)
	jobService.InitHandlers(cfg,logger)
	if err := jobService.Start(); err != nil {
		return nil, fmt.Errorf("failed to start job service: %w", err)
	}
	server := &Server{
		Config:		cfg,
		Logger:		logger,
		LoggerService: loggerService,
		DB:			db,
		Redis:		redisClient,
		Job:        jobService,
	}
	//start metric collection 
	//Runtime metrics are automatically collected by New Relic Go agent, so no additional setup is needed here for basic runtime metrics.
	return server, nil
}
