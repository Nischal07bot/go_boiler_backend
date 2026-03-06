package middlewares

//here we make a make a main middleware struct which will hold all the middlewares
import (
	"github.com/Nischal07bot/go_boiler_backend/internal/server"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Middleware struct {
	Global          *GlobalMiddlewares
	Auth            *AuthMiddleware
	ContextEnhancer *ContextEnhancer
	Tracing         *TracingMiddleware
	RateLimit       *RateLimitMiddleware
}

func NewMiddlewares(s *server.Server) *Middleware {
	var nrApp *newrelic.Application
	if s.LoggerService.GetApplication() != nil {
		nrApp = s.LoggerService.GetApplication()
	}
	return &Middleware{
		Global:          NewGlobalMiddlewares(s),
		Auth:            NewAuthMiddleware(s),
		ContextEnhancer: NewContextEnhancer(s),
		Tracing:         NewTracingMiddleware(s, nrApp),
		RateLimit:       NewRateLimitMiddleware(s),
	}
}
