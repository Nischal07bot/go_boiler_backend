package job

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/Nischal07bot/go_boiler_backend/internal/config"
)

type Jobservice struct {
	Client *asynq.Client
	server *asynq.Server
	Logger zerolog.Logger
}

func NewJobService(logger *zerolog.Logger,cfg *config.Config) *Jobservice {
	redisAddr := cfg.Redis.Address
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})
	server := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: redisAddr,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default": 3,
				"low": 1,
			},
		},
	)

	return &Jobservice{
		Client: client,
		server: server,
		Logger: *logger,
	}
}

func (j *Jobservice) Start() error{
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, j.handleWelcomeEmailTask)

	j.Logger.Info().Msg("Starting background job server")
	if err := j.server.Start(mux); err != nil {
		return err
	}

	return nil
}

func(j *Jobservice) Stop() {
	j.Logger.Info().Msg("Stopping background job server")
	j.server.Shutdown()
	j.Client.Close()
}