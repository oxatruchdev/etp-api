package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerDepartmentRoutes() {
	s.Echo.GET("/department", s.getDepartments)
	s.Echo.GET("/department/:id", s.getDepartment)
	s.Echo.POST("/department", s.createDepartment)
	s.Echo.PUT("/department/:id", s.updateDepartment)
}

func (s *Server) getDepartments(c echo.Context) error {
	var filter etp.DepartmentFilter
	if err := c.Bind(&filter); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	slog.Info("Getting departments")
	departments, n, err := s.DepartmentService.GetDepartments(c.Request().Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(200, echo.Map{"departments": departments, "count": n})
}

func (s *Server) createDepartment(c echo.Context) error {
	var department etp.Department

	if err := c.Bind(&department); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	if err := s.DepartmentService.CreateDepartment(c.Request().Context(), &department); err != nil {
		return Error(c, err)
	}
	return c.NoContent(http.StatusOK)
}

func (s *Server) getDepartment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid id"))
	}

	department, err := s.DepartmentService.GetDepartmentById(c.Request().Context(), id)
	if err != nil {
		return Error(c, err)
	}
	return c.JSON(200, department)
}

func (s *Server) updateDepartment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid id"))
	}

	var updDepartment etp.DepartmentUpdate

	if err := c.Bind(&updDepartment); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	department, err := s.DepartmentService.UpdateDepartment(c.Request().Context(), id, &updDepartment)
	if err != nil {
		return err
	}

	return c.JSON(200, department)
}
