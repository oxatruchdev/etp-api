package etpapi

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
	Departments   *[]Department    `json:"departments"`
	Professors    *[]Professor     `json:"professors"`
	SchoolRatings *[]SchoolRatings `json:"ratings"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SchoolFilter struct {
	ID        *int `json:"id"`
	CountryID *int `json:"countryId"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type SchoolService interface {
	GetSchoolById(ctx context.Context, id int) (*School, error)
	GetSchools(ctx context.Context, filter SchoolFilter) (*[]School, int, error)
	GetSchoolRatings(ctx context.Context, id int) (*[]SchoolRatings, error)
	CreateSchool(ctx context.Context, school *School) error
	UpdateSchool(ctx context.Context, id int, upd *SchoolUpdate) (*School, error)
	DeleteSchool(ctx context.Context, id int) error
}

type SchoolUpdate struct {
	Name         *string `json:"name"`
	Abbreviation *string `json:"abbreviation"`
}
