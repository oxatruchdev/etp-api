package main

import (
	"log/slog"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/Evalua-Tu-Profe/etp-api/db"
	"github.com/Evalua-Tu-Profe/etp-api/http"
)

type Main struct {
	DB         *db.DB
	HTTPServer *http.Server
	Config     *etp.Config
}

func NewMain() *Main {
	return &Main{
		DB:         db.NewDB(""),
		HTTPServer: http.NewServer(),
		Config:     etp.NewConfig(),
	}
}

func main() {
	slog.Info("Starting program")
	m := NewMain()

	slog.Info("Loading config")
	if err := m.Config.LoadConfig(); err != nil {
		slog.Error("Error loading config: ", "error", err)
		panic(err)
	}

	m.DB.DSN = m.Config.DatabaseURL

	slog.Info("Opening database")
	if error := m.DB.Open(); error != nil {
		slog.Error("Error opening database: ", "error", error)
		panic(error)
	}

	countryService := db.NewCountryService(m.DB)
	schoolService := db.NewSchoolService(m.DB)
	schoolRatingService := db.NewSchoolRatingService(m.DB)
	professorService := db.NewProfessorService(m.DB)
	professorServiceRating := db.NewProfessorRatingService(m.DB)
	departmentService := db.NewDepartmentService(m.DB)
	courseService := db.NewCourseService(m.DB)
	userService := db.NewUserService(m.DB)
	roleService := db.NewRoleService(m.DB)

	m.HTTPServer.CountryService = countryService
	m.HTTPServer.SchoolService = schoolService
	m.HTTPServer.SchoolRatingService = schoolRatingService
	m.HTTPServer.ProfessorService = professorService
	m.HTTPServer.ProfessorRatingService = professorServiceRating
	m.HTTPServer.DepartmentService = departmentService
	m.HTTPServer.CourseService = courseService
	m.HTTPServer.UserService = userService
	m.HTTPServer.RoleService = roleService

	m.HTTPServer.Config = m.Config

	slog.Info("Starting server")
	if err := m.HTTPServer.Start(m.Config.ServerPort); err != nil {
		slog.Error("Error starting server: ", "error", err)
		panic(err)
	}
}
