package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerCountryRoutes() {
	s.Echo.GET("/country", s.getCountries)
	s.Echo.GET("/country/:id", s.getCountry)
	s.Echo.POST("/country", s.createCountry)
	s.Echo.PUT("/country/:id", s.updateCountry)
}

func (s *Server) getCountries(c echo.Context) error {
	var filter etp.CountryFilter
	if err := c.Bind(&filter); err != nil {
		return err
	}
	countries, n, err := s.CountryService.GetCountries(c.Request().Context(), filter)
	slog.Info("countries", "countries", countries)
	if err != nil {
		slog.Error("error while getting countries", "error", err)
		return err
	}
	return c.JSON(200, echo.Map{"countries": countries, "count": n})
}

func (s *Server) getCountry(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	country, err := s.CountryService.GetCountryById(c.Request().Context(), id)
	if err != nil {
		slog.Error("error while getting countries", "error", err)
		return err
	}

	return c.JSON(200, country)
}

func (s *Server) createCountry(c echo.Context) error {
	var country etp.Country

	if err := c.Bind(&country); err != nil {
		return err
	}

	if err := s.CountryService.CreateCountry(c.Request().Context(), &country); err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func (s *Server) updateCountry(c echo.Context) error {
	var updCountry etp.CountryUpdate
	if err := c.Bind(&updCountry); err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	updatedCountry, err := s.CountryService.UpdateCountry(c.Request().Context(), id, &updCountry)
	if err != nil {
		return err
	}

	return c.JSON(200, updatedCountry)
}
