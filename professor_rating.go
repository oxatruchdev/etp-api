package etp

import (
	"context"
	"time"
)

type ProfessorRating struct {
	ID int `json:"id"`

	// Data related to the review itself
	WouldTakeAgain      bool   `json:"wouldTakeAgain" db:"would_take_again"`
	MandatoryAttendance bool   `json:"mandatoryAttendance" db:"mandatory_attendance"`
	Grade               string `json:"grade"`
	TextbookRequired    bool   `json:"textbookRequired" db:"textbook_required"`

	IsApproved     bool `json:"isApproved" db:"is_approved"`
	ApprovalsCount int  `json:"approvalsCount" db:"approvals_count"`
	UpdatedCount   int  `json:"updateCount" db:"updated_count"`

	// Relations
	Course      Course    `json:"course"`
	CourseId    int       `json:"courseId"`
	Professor   Professor `json:"professor"`
	ProfessorId int       `json:"professorId"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type ProfessorRatingFilter struct {
	ProfessorId *int `json:"professorId"`
	CourseId    *int `json:"courseId"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ProfessorRatingService interface {
	// Creates a new professor rating
	CreateProfessorRating(ctx context.Context, professorRating *ProfessorRating) error

	// Approves a professor rating
	// It is necessary to have at least 3 approvals in order to be approved
	ApproveProfessorRating(ctx context.Context, id int) error

	// Get all professor ratings
	// Can be filtered by the course and/or the professor
	GetProfessorRatings(ctx context.Context, filter ProfessorRatingFilter) ([]*ProfessorRating, int, error)

	// Deletes a professor rating
	DeleteProfessorRating(ctx context.Context, id int) error

	// Updates a professor rating
	// The rating will be put in a pending state until approved
	UpdateProfessorRating(ctx context.Context, id int, upd *ProfessorRatingUpdate) (*ProfessorRating, error)
}

type ProfessorRatingUpdate struct {
	WouldTakeAgain      *bool   `json:"wouldTakeAgain"`
	MandatoryAttendance *bool   `json:"mandatoryAttendance"`
	Grade               *string `json:"grade"`
	TextbookRequired    *bool   `json:"textbookRequired"`
}
