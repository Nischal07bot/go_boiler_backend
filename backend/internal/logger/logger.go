package logger

import (
	"fmt"
	"os"

	"github.com/Nischal07bot/go_boiler_backend/internal/config"
	"github.com/newrelic/go-agent/v3/newrelic"
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

	service.nrApp = app//struct initialization
	fmt.Printf("New Relic initialized for app: %s\n", cfg.ServiceName)
	return service //done
}
