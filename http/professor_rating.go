package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerProfessorRatingRoutes() {
	s.Echo.GET("/professor-rating", s.getProfessorRatings)
	s.Echo.POST("/professor-rating", s.createProfessorRating)
	s.Echo.PUT("/professor-rating/:id", s.updateProfessorRating)
	s.Echo.POST("/professor-rating/:id/approve", s.approveProfessorRating)
}

func (s *Server) getProfessorRatings(c echo.Context) error {
	var filter etp.ProfessorRatingFilter
	if err := c.Bind(&filter); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	slog.Info("Getting professor ratings", "filter", filter)

	professorRatings, n, err := s.ProfessorRatingService.GetProfessorRatings(c.Request().Context(), filter)
	if err != nil {
		return Error(c, err)
	}

	return c.JSON(200, echo.Map{"professorRatings": professorRatings, "count": n})
}

func (s *Server) createProfessorRating(c echo.Context) error {
	var professorRating etp.ProfessorRating
	if err := c.Bind(&professorRating); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	if err := s.ProfessorRatingService.CreateProfessorRating(c.Request().Context(), &professorRating); err != nil {
		return Error(c, err)
	}

	return c.NoContent(http.StatusCreated)
}

func (s *Server) updateProfessorRating(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid id"))
	}

	var professorRatingUpdate etp.ProfessorRatingUpdate
	if err := c.Bind(&professorRatingUpdate); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	if _, err := s.ProfessorRatingService.UpdateProfessorRating(c.Request().Context(), id, &professorRatingUpdate); err != nil {
		return Error(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (s *Server) approveProfessorRating(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid id"))
	}

	if err := s.ProfessorRatingService.ApproveProfessorRating(c.Request().Context(), id); err != nil {
		return Error(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
