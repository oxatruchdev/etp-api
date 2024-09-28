package etp

import (
	"context"
	"time"
)

type Professor struct {
	ID int `json:"id"`

	// Personal Information
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`

	School   *School `json:"school,omitempty" db:"-"`
	SchoolId int     `json:"schoolId"`

	// Relations
	Ratings     []*ProfessorRating `json:"ratings,omitempty" db:"-"`
	Departments []*Department      `json:"departments,omitempty" db:"-"`
	Courses     []*Course          `json:"courses,omitempty" db:"-"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type ProfessorFilter struct {
	ID *int `json:"id" param:"id"`

	Offset int `json:"offset" query:"offset"`
	Limit  int `json:"limit" query:"limit"`
}

type ProfessorService interface {
	// Gets a professor by id
	GetProfessorById(ctx context.Context, id int) (*Professor, error)

	// Gets all professors
	// Offset and Limit are used for pagination
	GetProfessors(ctx context.Context, filter ProfessorFilter) ([]*Professor, int, error)

	// Creates a professor
	CreateProfessor(ctx context.Context, professor *Professor) error

	// Updates a professor
	UpdateProfessor(ctx context.Context, id int, upd *ProfessorUpdate) (*Professor, error)

	// Deletes a professor
	DeleteProfessor(ctx context.Context, id int) error
}

type ProfessorUpdate struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
}
