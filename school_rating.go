package etp

import (
	"context"
	"time"
)

type SchoolRating struct {
	ID int `json:"id"`

	// Data related to the rating
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`

	// Data related to the school
	School   *School `json:"school"`
	SchoolID int     `json:"schoolId"`

	IsApproved bool `json:"isApproved" db:"is_approved"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type SchoolRatingFilter struct {
	SchoolID   *int  `json:"schoolId"`
	IsApproved *bool `json:"isApproved"`
	Offset     int   `json:"offset"`
	Limit      int   `json:"limit"`
}

type SchoolRatingService interface {
	CreateSchoolRating(ctx context.Context, schoolRatings *SchoolRating) error

	// Gets school ratings with pagination and filtering
	GetSchoolRatings(ctx context.Context, filter SchoolRatingFilter) ([]*SchoolRating, int, error)

	ApproveSchoolRating(ctx context.Context, id int) error
	UpdateSchoolRating(ctx context.Context, id int, upd *SchoolRatingUpdate) (*SchoolRatingUpdate, error)
	DeleteSchoolRating(ctx context.Context, id int) error
}

type SchoolRatingUpdate struct {
	Rating  *int    `json:"rating"`
	Comment *string `json:"comment"`
}
