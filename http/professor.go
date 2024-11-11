package http

import (
	"fmt"
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
	s.Mux.HandleFunc("POST /professor/{id}/review", s.CreateProfessorRating)
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
		IsApproved:  true,
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
		IsApproved:  true,
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
	if err != nil {
		slog.Info("error getting ratings", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func (s *Server) CreateProfessorRating(w http.ResponseWriter, r *http.Request) {
	// Getting professor
	slog.Info("Getting professor")
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var idInt int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validating form
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Validating form
	errors := make(map[string]string)

	// Validating rating
	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil {
		errors["rating"] = "Rating must be a number"
	}

	if rating < 1 || rating > 5 {
		errors["rating"] = "Rating must be between 1 and 5"
	}

	// Validating difficulty
	difficulty, err := strconv.Atoi(r.FormValue("difficulty"))
	if err != nil {
		errors["difficulty"] = "Difficulty must be a number"
	}

	fmt.Println("Difficulty", difficulty)

	if difficulty < 1 || difficulty > 5 {
		errors["difficulty"] = "Difficulty must be between 1 and 5"
	}

	// Validating would take again
	wouldTakeAgain := r.FormValue("wouldTakeAgain")
	if wouldTakeAgain != "true" && wouldTakeAgain != "false" {
		errors["wouldTakeAgain"] = "Would take again must be true or false"
	}

	// Validating course
	courseId, err := strconv.Atoi(r.FormValue("course"))
	if err != nil {
		errors["course"] = "Course must be a number"
	}

	// Validating comment
	comment := r.FormValue("comment")

	// Validating textbook required
	textbookRequired := r.FormValue("textBookRequired")
	if textbookRequired != "true" && textbookRequired != "false" {
		errors["textbookRequired"] = "Textbook required must be true or false"
	}

	// Validating mandatory attendance
	mandatoryAttendance := r.FormValue("mandatoryAttendance")
	if mandatoryAttendance != "true" && mandatoryAttendance != "false" {
		errors["mandatoryAttendance"] = "Mandatory attendance must be true or false"
	}

	// Handle multiple tags
	var tagIds []int
	formTags := r.Form["tags"] // This gets all tag values

	for _, tag := range formTags {
		tagId, err := strconv.Atoi(tag)
		if err != nil {
			errors["tags"] = "Tags must be numbers"
			break
		}
		tagIds = append(tagIds, tagId)
	}

	err = s.ProfessorRatingService.CreateProfessorRating(r.Context(), &etp.ProfessorRating{
		Rating:              rating,
		Comment:             comment,
		WouldTakeAgain:      wouldTakeAgain == "true",
		MandatoryAttendance: mandatoryAttendance == "true",
		Difficulty:          difficulty,
		TextbookRequired:    textbookRequired == "true",
		ProfessorId:         idInt,
		CourseId:            courseId,
	}, tagIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
