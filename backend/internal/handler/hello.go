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
