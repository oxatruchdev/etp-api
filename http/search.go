package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/Evalua-Tu-Profe/etp-api/cmd/web/partials"
	"github.com/a-h/templ"
)

type Search struct {
	Search string `form:"search"`
	Type   string `form:"type"`
}

func (s *Server) registerSearchRoutes() {
	s.Mux.HandleFunc("POST /search", s.search)
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	var search Search

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	if r.Form.Has("search") {
		search.Search = r.Form.Get("search")
	}

	if r.Form.Has("type") {
		search.Type = r.Form.Get("type")
	}

	if search.Type == "school" {
		slog.Info("Searching university")

		schools, _, err := s.SchoolService.GetSchools(r.Context(), etp.SchoolFilter{
			SchoolName: &search.Search,
		})
		if err != nil {
			slog.Error("Error while searching schools", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		slog.Info("Found schools", "schools", schools)

		// Getting countries for each school
		for i, school := range schools {
			country, err := s.CountryService.GetCountryById(r.Context(), school.CountryID)
			if err != nil {
				slog.Error("Error while searching schools", "error", err)
			}
			schools[i].Country = country
		}
		if err != nil {
			slog.Error("Error while searching schools", "error", err)
		}

		Render(w, r, http.StatusOK, partials.SchoolSearchResults(partials.SchoolSearchResultsProps{
			Results: formatSchoolsResults(schools),
		}))
		return
	}

	if search.Type == "professor" {
		professors, _, err := s.ProfessorService.GetProfessors(r.Context(), etp.ProfessorFilter{
			Name: &search.Search,
		})
		if err != nil {
			slog.Error("Error while searching professors", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		slog.Info("Found professors", "professors", professors)
		// getting school for each professor
		for i, professor := range professors {
			school, err := s.SchoolService.GetSchoolById(r.Context(), professor.SchoolId)
			slog.Info("Found school", "school", school)
			if err != nil {
				slog.Error("Error while searching schools", "error", err)
			}

			// getting country for each school
			country, err := s.CountryService.GetCountryById(r.Context(), school.CountryID)
			if err != nil {
				slog.Error("Error while searching schools", "error", err)
			}

			professors[i].School = school
			professors[i].School.Country = country
		}

		Render(w, r, http.StatusOK, partials.ProfessorSearchResults(partials.ProfessorSearchResultsProps{
			Results: formatProfessorResults(professors),
		}))
		return

	}
}

func formatSchoolsResults(schools []*etp.School) []partials.SchoolSearchResult {
	results := make([]partials.SchoolSearchResult, 0, len(schools))
	for _, school := range schools {
		results = append(results, partials.SchoolSearchResult{
			Name: school.Name,
			URL:  templ.SafeURL(fmt.Sprintf("/school/%d", school.ID)),
			Flag: fmt.Sprintf("fi-%s", school.Country.FlagCode),
		})
	}
	return results
}

func formatProfessorResults(professors []*etp.Professor) []partials.ProfessorSearchResult {
	results := make([]partials.ProfessorSearchResult, 0, len(professors))
	for _, professor := range professors {
		results = append(results, partials.ProfessorSearchResult{
			Name:       professor.FullName,
			URL:        templ.SafeURL(fmt.Sprintf("/professor/%d", professor.ID)),
			Flag:       fmt.Sprintf("fi-%s", professor.School.Country.FlagCode),
			University: professor.School.Name,
		})
	}
	return results
}
