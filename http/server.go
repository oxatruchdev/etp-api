package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/Evalua-Tu-Profe/etp-api/cmd/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Echo *echo.Echo

	// Services this server will use
	CountryService         etp.CountryService
	CourseService          etp.CourseService
	DepartmentService      etp.DepartmentService
	ProfessorService       etp.ProfessorService
	ProfessorRatingService etp.ProfessorRatingService
	SchoolService          etp.SchoolService
	SchoolRatingService    etp.SchoolRatingService
	UserService            etp.UserService
	RoleService            etp.RoleService

	Config *etp.Config
}

func NewServer() *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Caching static files
	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/css/*", echo.WrapHandler(fileServer))
	e.GET("/assets/*", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "max-age=31536000, public")
		fileServer.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	s := &Server{
		Echo: e,
	}

	slog.Info("Registering routes")
	{
		s.registerCountryRoutes()
		s.registerDepartmentRoutes()
		s.registerSchoolRoutes()
		s.registerCourseRoutes()
		s.registerProfessorRoutes()
		s.registerProfessorRatingRoutes()
		s.registerAuthRoutes()
		s.registerHomeRoutes()
		s.registerSearchRoutes()
	}

	return s
}

func (s *Server) Start(port int) error {
	return s.Echo.Start(fmt.Sprintf(":%d", port))
}
