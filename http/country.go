package http

import "github.com/labstack/echo/v4"

func (s *Server) registerCountryRoutes() {
	s.Echo.GET("/country", s.getCountries)
	s.Echo.GET("/country/:id", s.getCountry)
}

func (s *Server) getCountries(c echo.Context) error {
	return nil
}

func (s *Server) getCountry(c echo.Context) error {
	return nil
}
