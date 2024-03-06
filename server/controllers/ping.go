package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PingController struct{}

func NewPingController() *PingController {
	return &PingController{}
}

func (p *PingController) GetGroup() string {
	return "api/ping"
}
func (p *PingController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{
		&Handler{
			Method:  "GET",
			Path:    "",
			Handler: p.ping,
		},
	}
}
func (p *PingController) GetMiddleware() []echo.MiddlewareFunc {
	return nil
}

func (p *PingController) ping(e echo.Context) error {
	return e.JSON(http.StatusOK, "pong")
}
