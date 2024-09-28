package http

import (
	"net/http"
	"strconv"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerProfessorRoutes() {
	s.Echo.GET("/professor", s.getProfessors)
	s.Echo.GET("/professor/:id", s.getProfessor)
	s.Echo.POST("/professor", s.createProfessor)
	s.Echo.PUT("/professor/:id", s.updateProfessor)
}

func (s *Server) getProfessors(c echo.Context) error {
	var filter etp.ProfessorFilter
	if err := c.Bind(&filter); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	professors, n, err := s.ProfessorService.GetProfessors(c.Request().Context(), filter)
	if err != nil {
		return Error(c, err)
	}

	return c.JSON(200, echo.Map{"professors": professors, "count": n})
}

func (s *Server) getProfessor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid id"))
	}

	professor, err := s.ProfessorService.GetProfessorById(c.Request().Context(), id)
	if err != nil {
		return Error(c, err)
	}

	return c.JSON(200, professor)
}

func (s *Server) createProfessor(c echo.Context) error {
	var professor etp.Professor
	if err := c.Bind(&professor); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	if err := s.ProfessorService.CreateProfessor(c.Request().Context(), &professor); err != nil {
		return Error(c, err)
	}

	return c.NoContent(http.StatusCreated)
}

func (s *Server) updateProfessor(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid id"))
	}

	var upd etp.ProfessorUpdate
	if err := c.Bind(&upd); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	professor, err := s.ProfessorService.UpdateProfessor(c.Request().Context(), id, &upd)
	if err != nil {
		return Error(c, err)
	}

	return c.JSON(http.StatusOK, professor)
}
