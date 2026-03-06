package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HelloHandler struct {
	Handler
}

func NewHelloHandler(h Handler) *HelloHandler {
	return &HelloHandler{Handler: h}
}

func (h *HelloHandler) Hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello my friend")
}
/*
┌─────────────────────────────────────────────────────────────────────────────┐
│                           APPLICATION STARTUP                                │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  main()                                                                      │
│  ├── config.Load()                    → Load environment config              │
│  ├── logger.NewLoggerService()        → Setup logging                        │
│  ├── server.New()                     → Create Server (DB, Redis, Jobs)      │
│  ├── service.NewServices()            → Create business logic layer          │
│  ├── handler.NewHandlers() ─────────────────────────────────────────────┐    │
│  ├── router.NewRouter() ────────────────────────────────────────────┐   │    │
│  └── router.Start()                   → Start HTTP server           │   │    │
└──────────────────────────────────────────────────────────────────────│───│────┘
                                                                       │   │
                    ┌──────────────────────────────────────────────────┘   │
                    ▼                                                      │
┌─────────────────────────────────────────────────────────────────────┐    │
│  handler.NewHandlers(server, services)                              │    │
│  ├── NewHandler(server) ──────────────────────────────────────┐     │    │
│  │                                                            │     │    │
│  │   ┌──────────────────────────────────────────────────────┐ │     │    │
│  │   │  NewHandler(s)                                       │ │     │    │
│  │   │  └── return Handler{server: s}                       │ │     │    │
│  │   │      (Base handler with server access)               │ │     │    │
│  │   └──────────────────────────────────────────────────────┘ │     │    │
│  │                                                            │     │    │
│  └── NewHelloHandler(h) ──────────────────────────────────────│─┐   │    │
│                                                               │ │   │    │
│      ┌────────────────────────────────────────────────────────│─│───┘    │
│      ▼                                                        │ │        │
│  ┌──────────────────────────────────────────────────────────┐ │ │        │
│  │  NewHelloHandler(h Handler)                              │ │ │        │
│  │  └── return &HelloHandler{Handler: h}                    │◄┘ │        │
│  │      (HelloHandler embeds base Handler)                  │   │        │
│  └──────────────────────────────────────────────────────────┘   │        │
│                                                                 │        │
│  └── return &Handlers{Hello: helloHandler}                      │        │
└─────────────────────────────────────────────────────────────────┘        │
                                                                           │
                    ┌──────────────────────────────────────────────────────┘
                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  router.NewRouter(server, handlers, services)                                │
│  ├── middlewares.NewMiddlewares(s)    → Setup middlewares                    │
│  ├── echo.New()                       → Create Echo router                   │
│  ├── router.Use(...)                  → Apply global middlewares             │
│  └── registerSystemRoutes(router, handlers) ──────────────────────────┐      │
└───────────────────────────────────────────────────────────────────────│──────┘
                                                                        │
                    ┌───────────────────────────────────────────────────┘
                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  registerSystemRoutes(r, h)                                                  │
│  └── r.GET("/hello", h.Hello.Hello)                                          │
│          │       │         │                                                 │
│          │       │         └── HelloHandler.Hello method                     │
│          │       └── HelloHandler instance                                   │
│          └── Handlers struct                                                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                           HTTP REQUEST TIME                                  │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│  GET /hello                                                                  │
│      │                                                                       │
│      ▼                                                                       │
│  ┌───────────────────────────────────────┐                                   │
│  │  Middlewares execute in order:        │                                   │
│  │  1. RateLimiter                       │                                   │
│  │  2. CORS                              │                                   │
│  │  3. Secure                            │                                   │
│  │  4. RequestID                         │                                   │
│  │  5. NewRelicMiddleware                │                                   │
│  │  6. EnhanceTracing                    │                                   │
│  │  7. EnhanceContext                    │                                   │
│  │  8. RequestLogger                     │                                   │
│  │  9. Recover                           │                                   │
│  └───────────────────────────────────────┘                                   │
│      │                                                                       │
│      ▼                                                                       │
│  ┌───────────────────────────────────────┐                                   │
│  │  HelloHandler.Hello(c echo.Context)   │                                   │
│  │  └── return c.String(200,             │                                   │
│  │         "hello my friend")            │                                   │
│  └───────────────────────────────────────┘                                   │
│      │                                                                       │
│      ▼                                                                       │
│  Response: "hello my friend" (HTTP 200)                                      │
└─────────────────────────────────────────────────────────────────────────────┘
*/