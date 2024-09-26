package http

import (
	"fmt"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/labstack/echo/v4"
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
}

func NewServer() *Server {
	e := echo.New()

	s := &Server{
		Echo: e,
	}

	{
		s.registerCountryRoutes()
	}

	return s
}

func (s *Server) Start(port int) error {
	return s.Echo.Start(fmt.Sprintf(":%d", port))
}
