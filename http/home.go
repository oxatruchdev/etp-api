package http

import (
	"log/slog"
	"net/http"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/Evalua-Tu-Profe/etp-api/cmd/web"
)

func (s *Server) registerHomeRoutes() {
	s.Mux.HandleFunc("GET /", s.home)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	slog.Info("Hitting home route")

	schools, _, err := s.SchoolService.GetSchools(r.Context(), etp.SchoolFilter{
		Offset: 0,
		Limit:  5,
	})
	if err != nil {
		slog.Error("Error while searching schools", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	slog.Info("Found schools", "schools", schools)

	professorCountBySchools := make(map[int]int)

	for _, school := range schools {
		professorCountBySchools[school.ID], err = s.SchoolService.GetSchoolProfessorsCount(r.Context(), school.ID)
		if err != nil {
			slog.Error("Error while searching schools", "error", err)
			return
		}
	}

	// Get latest reviews
	ratings, err := s.ProfessorRatingService.GetLatestProfessorsRatings(r.Context(), etp.ProfessorRatingFilter{
		Offset: 0,
		Limit:  5,
	})
	if err != nil {
		slog.Error("Error while getting latest reviews", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	Render(w, r, http.StatusOK, web.Home(web.HomeProps{
		Schools:                schools,
		ProfessorCountBySchool: professorCountBySchools,
		Ratings:                ratings,
	}))
}
