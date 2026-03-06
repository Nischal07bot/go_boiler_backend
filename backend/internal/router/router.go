package router

import (
	"net/http"

	"github.com/Nischal07bot/go_boiler_backend/internal/handler"
	"github.com/Nischal07bot/go_boiler_backend/internal/middlewares"
	"github.com/Nischal07bot/go_boiler_backend/internal/server"
	"github.com/Nischal07bot/go_boiler_backend/internal/service"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func NewRouter(s *server.Server, h *handler.Handlers, services *service.Services) *echo.Echo {
	mw := middlewares.NewMiddlewares(s)

	router := echo.New()

	router.HTTPErrorHandler = mw.Global.GlobalErrorHandler

	// global middlewares
	router.Use(
		echoMiddleware.RateLimiterWithConfig(echoMiddleware.RateLimiterConfig{
			Store: echoMiddleware.NewRateLimiterMemoryStore(rate.Limit(20)),
			DenyHandler: func(c echo.Context, identifier string, err error) error {
				// Record rate limit hit metrics
				if rateLimitMiddleware := mw.RateLimit; rateLimitMiddleware != nil {
					rateLimitMiddleware.RecordRateLimitHit(c.Path())
				}

				s.Logger.Warn().
					Str("request_id", middlewares.GetRequestID(c)).
					Str("identifier", identifier).
					Str("path", c.Path()).
					Str("method", c.Request().Method).
					Str("ip", c.RealIP()).
					Msg("rate limit exceeded")

				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded")
			},
		}),
		mw.Global.CORS(),
		mw.Global.Secure(),
		middlewares.RequestID(),
		mw.Tracing.NewRelicMiddleware(),
		mw.Tracing.EnhanceTracing(),
		mw.ContextEnhancer.EnhanceContext(),
		mw.Global.RequestLogger(),
		mw.Global.Recover(),
	)

	// register system routes
	registerSystemRoutes(router, h)

	// register versioned routes
	router.Group("/api/v1")

	return router
}
