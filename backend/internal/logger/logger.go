package logger

import (
	"fmt"
	"io"
	"os"
	"time"
	"strings"

	"github.com/Nischal07bot/go_boiler_backend/internal/config"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type LoggerService struct {
	nrApp *newrelic.Application
}

// this creates a new logger service with the newrelic intergrayion
func NewLoggerService(cfg *config.ObservabilityConfig) *LoggerService {
	service := &LoggerService{}
	if cfg.NewRelic.LicenseKey == "" {
		fmt.Println("New Relic license key not provided, skipping New Relic integration")
		return service
	}
	var configOptions []newrelic.ConfigOption
	configOptions = append(configOptions,
		newrelic.ConfigAppName(cfg.ServiceName),
		newrelic.ConfigLicense(cfg.NewRelic.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(cfg.NewRelic.DistributedTracingEnabled),
		newrelic.ConfigAppLogForwardingEnabled(cfg.NewRelic.AppLogForwardingEnabled),
	)
	/*The application should use LoggerService instead of New Relic directly to avoid tight coupling with a third-party library, improve testability, centralize logging behavior, and
	allow observability tools to be swapped or extended without changing business logic.*/

	if cfg.NewRelic.DebugLogging {
		configOptions = append(configOptions, newrelic.ConfigDebugLogger(os.Stdout))
	}

	app, err := newrelic.NewApplication(configOptions...)
	if err != nil {
		fmt.Printf("Failed to intialize New Relic %v\n", err)
		return service
	}

	service.nrApp = app //struct initialization
	fmt.Printf("New Relic initialized for app: %s\n", cfg.ServiceName)
	return service //done
}

// this is used to gracefully shutdown the newrelic application when the server is shutting down to ensure that all pending logs and traces are sent to newrelic before the application exits.
func (ls *LoggerService) Shutdown() {
	if ls.nrApp != nil {
		ls.nrApp.Shutdown(10 * time.Second)
	}
}

// this returns the newrelic application instance to be used in other packages for logging and tracing purposes without directly coupling those packages to the newrelic library. This allows for better modularity and testability of the code.
func (ls *LoggerService) GetApplication() *newrelic.Application {
	return ls.nrApp //this is kind of a getter since nrApp cant be exposed not in capital letter

}

func NewLogger(level string, isProd bool) zerolog.Logger {
	return NewLoggerWithService(&config.ObservabilityConfig{
		Logging: config.LoggingConfig{
			Level: level,
		},
		Environment: func() string {
			if isProd {
				return "production"
			}
			return "development"
		}(),
	}, nil)
}

func NewLoggerWithConfig(cfg *config.ObservabilityConfig) zerolog.Logger {
	return NewLoggerWithService(cfg, nil)
}

func NewLoggerWithService(cfg *config.ObservabilityConfig, ls *LoggerService) zerolog.Logger {
	var loglevel zerolog.Level
	level := cfg.Logging.Level

	switch level {
	case "debug":
		loglevel = zerolog.DebugLevel
	case "info":
		loglevel = zerolog.InfoLevel
	case "warn":
		loglevel = zerolog.WarnLevel
	case "error":
		loglevel = zerolog.ErrorLevel
	default:
		loglevel = zerolog.InfoLevel//storing zerlog.Inolevel since the variable loglevel
		//will be further used by zerolog hence cant use string level since zerolog needs the log level in its own type
	}

	//dont set global log level since we want each logger its own level
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var writer io.Writer

	//setup base writer
	var baseWriter io.Writer
	if cfg.IsProduction() && cfg.Logging.Format == "json" {
		baseWriter = os.Stdout

		//wrap with newrelic log forwarder if enabled
		if ls != nil && ls.nrApp != nil {
			nrWriter := zerologWriter.New(baseWriter, ls.nrApp)
			writer = nrWriter
		} else {
			writer = baseWriter
		}

	} else {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
		writer = consoleWriter
	}

	logger := zerolog.New(writer).Level(loglevel).With().Timestamp().Str("service", cfg.ServiceName).
		Str("environment", cfg.Environment).Logger()

	if !cfg.IsProduction() {
		logger = logger.With().Stack().Logger()
	}
	return logger

}


func WithTraceContext(logger zerolog.Logger, ctx *newrelic.Transaction) zerolog.Logger {
	if ctx == nil {
		return logger
	}
	//get trace metadata from the transaction 

	metadata := ctx.GetTraceMetadata()

	return logger.With().Str("trace.id", metadata.TraceID).Str("span.id", metadata.SpanID).Logger()
}

func FormatSQLWithArgs(sql string,args []any) string {
	result := sql
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		value := fmt.Sprintf("'%v", arg)
		result = strings.Replace(result,placeholder,value,1)
	}

	return result
}
