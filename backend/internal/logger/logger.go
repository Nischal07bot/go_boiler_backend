package logger

import (
	"fmt"

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
	return service //will implement this later 
}
