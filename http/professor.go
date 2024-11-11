package http

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/Evalua-Tu-Profe/etp-api/cmd/web"
	"github.com/Evalua-Tu-Profe/etp-api/cmd/web/components"
)

func (s *Server) registerProfessorRoutes() {
	s.Mux.HandleFunc("GET /professor/{id}", s.getProfessor)
	s.Mux.HandleFunc("GET /professor/{id}/reviews", s.getProfessorReviews)
	s.Mux.HandleFunc("GET /professor/{id}/add-review", s.HandleAddProfessorReview)
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

	// Getting courses for that professor
	courses, err := s.ProfessorService.GetProfessorCourses(r.Context(), idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	professor.Courses = courses

	// Getting professor's school
	school, err := s.SchoolService.GetSchoolById(r.Context(), professor.SchoolId)
	if err != nil {
		slog.Info("error getting school", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	professor.School = school

	// Getting professor ratings
	ratings, err := s.ProfessorRatingService.GetProfessorRatingsWithStats(r.Context(), etp.ProfessorRatingFilter{
		ProfessorId: &idInt,
	})
	if err != nil {
		slog.Info("error getting ratings", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	professor.Ratings = ratings.Ratings

	// Getting professor's most popular tags
	tags, err := s.ProfessorService.GetProfessorTags(r.Context(), idInt)
	if err != nil {
		slog.Info("error getting tags", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	professor.PopularTags = tags

	props := web.ProfessorPageProps{
		Professor:        professor,
		School:           professor.School,
		RatingsWithStats: *ratings,
	}
	Render(w, r, http.StatusOK, web.ProfessorPage(props))
}

func (s *Server) getProfessorReviews(w http.ResponseWriter, r *http.Request) {
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

	filter := etp.ProfessorRatingFilter{
		ProfessorId: &idInt,
	}

	course := r.URL.Query().Get("course")
	if course != "" {
		courseId, err := strconv.Atoi(course)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		filter.CourseId = &courseId
	}

	// Getting professor ratings
	ratings, err := s.ProfessorRatingService.GetProfessorRatingsWithStats(r.Context(), filter)

	Render(w, r, 200, components.RatingsList(ratings.Ratings, ratings.TotalCount))
}

func (s *Server) HandleAddProfessorReview(w http.ResponseWriter, r *http.Request) {
	// Getting professor
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

	// Getting professor's school
	school, err := s.SchoolService.GetSchoolById(r.Context(), professor.SchoolId)
	if err != nil {
		slog.Info("error getting school", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	professor.School = school

	// Getting courses for the professor's department
	courses, _, err := s.CourseService.GetCourses(r.Context(), etp.CourseFilter{
		DepartmentId: &professor.Department.ID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	professor.Courses = courses

	// Getting tags
	tags, err := s.TagService.GetTags(r.Context())
	if err != nil {
		slog.Info("error getting tags", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	Render(w, r, 200, web.AddProfessorReviewPage(web.AddProfessorReviewPageProps{
		Professor: professor,
		Tags:      tags,
	}))
}
