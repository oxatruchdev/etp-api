package http

func (s *Server) registerCourseRoutes() {
	// s.Echo.GET("/course", s.getCourses)
	// s.Echo.GET("/course/:id", s.getCourse)
	// s.Echo.POST("/course", s.createCourse)
	// s.Echo.PUT("/course/:id", s.updateCourse)
}

// func (s *Server) getCourses(c echo.Context) error {
// 	var filter etp.CourseFilter
// 	if err := c.Bind(&filter); err != nil {
// 		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
// 	}
//
// 	courses, n, err := s.CourseService.GetCourses(c.Request().Context(), filter)
// 	if err != nil {
// 		return Error(c, err)
// 	}
//
// 	return c.JSON(200, echo.Map{"courses": courses, "count": n})
// }
//
// func (s *Server) getCourse(c echo.Context) error {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return Error(c, etp.Errorf(etp.EINVALID, "invalid id"))
// 	}
//
// 	course, err := s.CourseService.GetCourseById(c.Request().Context(), id)
// 	if err != nil {
// 		return Error(c, err)
// 	}
//
// 	return c.JSON(200, course)
// }
//
// func (s *Server) createCourse(c echo.Context) error {
// 	var course etp.Course
// 	if err := c.Bind(&course); err != nil {
// 		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
// 	}
//
// 	if err := s.CourseService.CreateCourse(c.Request().Context(), &course); err != nil {
// 		return Error(c, err)
// 	}
//
// 	return c.NoContent(http.StatusCreated)
// }
//
// func (s *Server) updateCourse(c echo.Context) error {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		return Error(c, etp.Errorf(etp.EINVALID, "invalid id"))
// 	}
//
// 	var course *etp.CourseUpdate
// 	if err := c.Bind(&course); err != nil {
// 		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
// 	}
//
// 	courseUpdated, err := s.CourseService.UpdateCourse(c.Request().Context(), id, course)
// 	if err != nil {
// 		return Error(c, err)
// 	}
//
// 	return c.JSON(200, courseUpdated)
// }
