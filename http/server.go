package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/Evalua-Tu-Profe/etp-api/cmd/web"
	"github.com/Evalua-Tu-Profe/etp-api/middleware"
)

type Server struct {
	Mux    *http.ServeMux
	Server *http.Server

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
	s := &Server{
		Mux:    http.NewServeMux(),
		Server: &http.Server{},
		Config: etp.NewConfig(),
	}

	// Registering static assets
	fileServer := http.FileServer(http.FS(web.Files))
	s.Mux.Handle("GET /assets/", fileServer)

	{
		s.registerAuthRoutes()
		s.registerHomeRoutes()
		s.registerSearchRoutes()
	}

	slog.Info("Creating middleware stack")
	stack := middleware.CreateStack(middleware.Logging, middleware.AuthOptional)

	s.Server = &http.Server{
		Addr:           ":8080", // You can update this dynamically in `Start`
		Handler:        stack(s.Mux),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return s
}

func (s *Server) Start(port int) error {
	// Update the address dynamically based on the port provided
	s.Server.Addr = fmt.Sprintf(":%d", port)

	fmt.Printf("Starting server on port %d...\n", port)
	return s.Server.ListenAndServe()
}
