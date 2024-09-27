package http

import (
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
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	schools, n, err := s.SchoolService.GetSchools(c.Request().Context(), filter)
	if err != nil {
		return Error(c, err)
	}

	return c.JSON(200, echo.Map{"schools": schools, "count": n})
}

func (s *Server) getSchool(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid id"))
	}

	school, err := s.SchoolService.GetSchoolById(c.Request().Context(), id)
	if err != nil {
		return Error(c, err)
	}

	return c.JSON(200, school)
}

func (s *Server) createSchool(c echo.Context) error {
	var school etp.School
	if err := c.Bind(&school); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	if err := s.SchoolService.CreateSchool(c.Request().Context(), &school); err != nil {
		return Error(c, err)
	}

	return c.NoContent(http.StatusCreated)
}

func (s *Server) updateSchool(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid id"))
	}

	var upd etp.SchoolUpdate
	if err := c.Bind(&upd); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	school, err := s.SchoolService.UpdateSchool(c.Request().Context(), id, &upd)
	if err != nil {
		return Error(c, err)
	}

	return c.JSON(200, school)
}
