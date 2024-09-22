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

	IsApproved bool `json:"isApproved"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SchoolRatingService interface {
	GetSchoolRatings(ctx context.Context, id int) ([]*SchoolRating, error)
	CreateSchoolRating(ctx context.Context, schoolRatings *SchoolRating) error
	ApproveSchoolRating(ctx context.Context, id int) error
	UpdateSchoolRating(ctx context.Context, id int, upd *SchoolRatingUpdate) (*SchoolRatingUpdate, error)
	DeleteSchoolRating(ctx context.Context, id int) error
}

type SchoolRatingUpdate struct {
	Rating  *int    `json:"rating"`
	Comment *string `json:"comment"`
}
