package etpapi

import (
	"context"
	"time"
)

type SchoolRatings struct {
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

type SchoolRatingsService interface {
	GetSchoolRatings(ctx context.Context, id int) (*[]SchoolRatings, error)
	CreateSchoolRatings(ctx context.Context, schoolRatings *SchoolRatings) error
	ApproveSchoolRating(ctx context.Context, id int) error
	UpdateSchoolRatings(ctx context.Context, id int, upd *SchoolRatingsUpdate) (*SchoolRatings, error)
	DeleteSchoolRatings(ctx context.Context, id int) error
}

type SchoolRatingsUpdate struct {
	Rating  *int    `json:"rating"`
	Comment *string `json:"comment"`
}
