package http

import (
	"log/slog"
	"net/http"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/Evalua-Tu-Profe/etp-api/cmd/web/partials"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Search struct {
	Search string `form:"search"`
	Type   string `form:"type"`
}

func (s *Server) registerSearchRoutes() {
	s.Echo.POST("/search", s.search)
}

func (s *Server) search(c echo.Context) error {
	var search Search

	if err := c.Bind(&search); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	slog.Info("Searching", "search", search)

	results := []struct {
		Name string
		URL  templ.SafeURL
	}{
		{Name: "Departments", URL: templ.URL("/departments")},
		{Name: "Courses", URL: templ.URL("/courses")},
		{Name: "Professors", URL: templ.URL("/professors")},
		{Name: "Schools", URL: templ.URL("/schools")},
		{Name: "Ratings", URL: templ.URL("/ratings")},
	}

	return Render(c, http.StatusOK, partials.SearchResults(partials.SearchResultsProps{
		Results: results,
	}))
}
