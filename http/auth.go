package http

import (
	"net/http"

	"github.com/Evalua-Tu-Profe/etp-api/cmd/web"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerAuthRoutes() {
	s.Echo.GET("/register", s.login)
	s.Echo.POST("/register", s.register)

	s.Echo.GET("/login", s.login)
	s.Echo.POST("/login", s.login)
}

func (s *Server) login(c echo.Context) error {
	return Render(c, http.StatusOK, web.Register())
}

func (s *Server) register(c echo.Context) error {
	return Render(c, http.StatusOK, web.Register())
}
