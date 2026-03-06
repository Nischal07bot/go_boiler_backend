package handler

import (
	"github.com/Nischal07bot/go_boiler_backend/internal/server"
	"github.com/Nischal07bot/go_boiler_backend/internal/service"
)

type Handlers struct {
	Hello *HelloHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	h := NewHandler(s)
	return &Handlers{
		Hello: NewHelloHandler(h),
	}
}
