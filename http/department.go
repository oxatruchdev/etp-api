package http

import (
	"net/http"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerDepartmentRoutes() {
	s.Echo.GET("/department", s.getDepartments)
	s.Echo.GET("/department/:id", s.getDepartment)
	s.Echo.POST("/department", s.createDepartment)
}

func (s *Server) getDepartments(c echo.Context) error {
	departments, n, err := s.DepartmentService.GetDepartments(c.Request().Context(), etp.DepartmentFilter{})
	if err != nil {
		return err
	}

	return c.JSON(200, echo.Map{"departments": departments, "count": n})
}

func (s *Server) createDepartment(c echo.Context) error {
	var department etp.Department

	if err := c.Bind(&department); err != nil {
		return err
	}

	if err := s.DepartmentService.CreateDepartment(c.Request().Context(), &department); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (s *Server) getDepartment(c echo.Context) error {
	return nil
}
