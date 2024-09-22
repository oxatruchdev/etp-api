package etp

import (
	"context"
	"time"
)

type ProfessorRating struct {
	ID int `json:"id"`

	// Data related to the review itself
	WouldTakeAgain      bool   `json:"wouldTakeAgain"`
	MandatoryAttendance bool   `json:"mandatoryAttendance"`
	Grade               string `json:"grade"`
	TextbookRequired    bool   `json:"textbookRequired"`

	IsApproved     bool `json:"isApproved"`
	ApprovalsCount int  `json:"approvalsCount"`

	// Relations
	Course   Course `json:"course"`
	CourseId int    `json:"courseId"`

	Professor   Professor `json:"professor"`
	ProfessorId int       `json:"professorId"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
	// Can be filter by the course and/or the professor
	GetProfessorRatings(ctx context.Context, filter ProfessorFilter) ([]*ProfessorRating, error)

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
