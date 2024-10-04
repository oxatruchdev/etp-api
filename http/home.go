package http

import (
	"net/http"

	"github.com/Evalua-Tu-Profe/etp-api/cmd/web"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerHomeRoutes() {
	s.Echo.GET("/", s.home)
}

func (s *Server) home(c echo.Context) error {
	return Render(c, http.StatusOK, web.Home())
}
