package http

import (
	"log/slog"
	"net/http"

	"github.com/Evalua-Tu-Profe/etp-api/cmd/web"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerHomeRoutes() {
	s.Echo.GET("/", s.home, s.AuthMiddleware)
}

func (s *Server) home(c echo.Context) error {
	isAuthenticated := c.Get("isAuthenticated").(bool)
	slog.Info("Rendering home", "isAuth", isAuthenticated)
	return Render(c, http.StatusOK, web.Home(web.HomeProps{
		IsAuthenticated: isAuthenticated,
	}))
}
