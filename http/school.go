package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerSchoolRoutes() {
	s.Echo.GET("/school", s.getSchools)
	s.Echo.GET("/school/:id", s.getSchool)
	s.Echo.POST("/school", s.createSchool)
	s.Echo.PUT("/school/:id", s.updateSchool)
}

func (s *Server) getSchools(c echo.Context) error {
	var filter etp.SchoolFilter
	if err := c.Bind(&filter); err != nil {
		return err
	}
	schools, n, err := s.SchoolService.GetSchools(c.Request().Context(), filter)
	slog.Info("schools", "schools", schools)
	if err != nil {
		slog.Error("error while getting schools", "error", err)
		return err
	}
	return c.JSON(200, echo.Map{"schools": schools, "count": n})
}

func (s *Server) getSchool(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	school, err := s.SchoolService.GetSchoolById(c.Request().Context(), id)
	if err != nil {
		slog.Error("error while getting schools", "error", err)
		return err
	}

	return c.JSON(200, school)
}

func (s *Server) createSchool(c echo.Context) error {
	var school etp.School

	if err := c.Bind(&school); err != nil {
		return err
	}

	if err := s.SchoolService.CreateSchool(c.Request().Context(), &school); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (s *Server) updateSchool(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	var upd etp.SchoolUpdate
	if err := c.Bind(&upd); err != nil {
		return err
	}

	school, err := s.SchoolService.UpdateSchool(c.Request().Context(), id, &upd)
	if err != nil {
		return err
	}

	return c.JSON(200, school)
}
