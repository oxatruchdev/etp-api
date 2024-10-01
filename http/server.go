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
}

func NewServer() *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/assets/*", echo.WrapHandler(fileServer))

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
	}

	return s
}

func (s *Server) Start(port int) error {
	return s.Echo.Start(fmt.Sprintf(":%d", port))
}
