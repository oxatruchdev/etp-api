package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/Evalua-Tu-Profe/etp-api/cmd/web"
)

func (s *Server) registerProfessorRoutes() {
	s.Mux.HandleFunc("GET /professor/{id}", s.getProfessor)
}

func (s *Server) getProfessor(w http.ResponseWriter, r *http.Request) {
	slog.Info("Getting professor")
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("Get professor", "id", id)

	professor, err := s.ProfessorService.GetProfessorById(r.Context(), idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Got professor", "id", id, "professor", professor)

	slog.Info("Getting school")
	// Getting professor's school
	school, err := s.SchoolService.GetSchoolById(r.Context(), professor.SchoolId)
	if err != nil {
		slog.Info("error getting school", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	professor.School = school

	slog.Info("Got school", "id", professor.School.ID, "school", professor.School)

	// Getting professor ratings
	ratings, err := s.ProfessorRatingService.GetProfessorRatingsWithStats(r.Context(), etp.ProfessorRatingFilter{
		ProfessorId: &idInt,
	})

	slog.Info("Got ratings", "ratings", ratings)
	if err != nil {
		slog.Info("error getting ratings", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	professor.Ratings = ratings.Ratings

	props := web.ProfessorPageProps{
		Professor:        professor,
		School:           professor.School,
		RatingsWithStats: *ratings,
	}
	Render(w, r, http.StatusOK, web.ProfessorPage(props))
}
