package etpapi

import (
	"context"
	"time"
)

type Professor struct {
	ID int `json:"id"`

	// Personal Information
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	School   *School `json:"school"`
	SchoolId int     `json:"schoolId"`

	// Relations
	Ratings *[]ProfessorRating `json:"ratings"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProfessorFilter struct {
	ID *int `json:"id"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ProfessorService interface {
	// Gets a professor by id
	GetProfessorById(ctx context.Context, id int) (*Professor, error)

	// Gets all professors
	// Offset and Limit are used for pagination
	GetProfessors(ctx context.Context, filter ProfessorFilter) (*[]Professor, int, error)
	GetProfessorRatings(ctx context.Context, id int) (*[]ProfessorRating, error)
	CreateProfessor(ctx context.Context, professor *Professor) error
	UpdateProfessor(ctx context.Context, id int, upd *ProfessorUpdate) (*Professor, error)
	DeleteProfessor(ctx context.Context, id int) error
}

type ProfessorUpdate struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}
