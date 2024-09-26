package http

import (
	"log/slog"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerCountryRoutes() {
	s.Echo.GET("/country", s.getCountries)
	s.Echo.GET("/country/:id", s.getCountry)
	s.Echo.POST("/country", s.createCountry)
}

func (s *Server) getCountries(c echo.Context) error {
	countries, _, err := s.CountryService.GetCountries(c.Request().Context(), etp.CountryFilter{
		Offset: 0,
		Limit:  10,
	})
	slog.Info("countries", countries)
	if err != nil {
		slog.Error("error while getting countries", "error", err)
		return err
	}
	return c.JSON(200, countries)
}

func (s *Server) getCountry(c echo.Context) error {
	return nil
}

func (s *Server) createCountry(c echo.Context) error {
	var country etp.Country

	if err := c.Bind(&country); err != nil {
		return err
	}

	if err := s.CountryService.CreateCountry(c.Request().Context(), &country); err != nil {
		return err
	}

	return c.JSON(200, country)
}
