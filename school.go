package etp

import (
	"context"
	"time"
)

type School struct {
	ID int `json:"id"`

	// Name and abbreviation are used for display
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`

	// Country school belongs to
	Country   *Country `json:"country"`
	CountryID int      `json:"countryId"`

	// Relations
	Departments   []*Department   `json:"departments"`
	Professors    []*Professor    `json:"professors"`
	SchoolRatings []*SchoolRating `json:"ratings"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type SchoolFilter struct {
	SchoolId  *int `json:"id"`
	CountryID *int `json:"countryId"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type SchoolService interface {
	// Gets a school by id
	GetSchoolById(ctx context.Context, id int) (*School, error)

	// Gets all schools
	// Offset and Limit are used for pagination
	GetSchools(ctx context.Context, filter SchoolFilter) ([]*School, int, error)

	// Creates a school
	CreateSchool(ctx context.Context, school *School) error

	// Updates a school
	UpdateSchool(ctx context.Context, id int, upd *SchoolUpdate) (*School, error)

	// Deletes a school
	DeleteSchool(ctx context.Context, id int) error
}

type SchoolUpdate struct {
	Name         *string `json:"name"`
	Abbreviation *string `json:"abbreviation"`
}
