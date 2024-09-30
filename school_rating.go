package etp

import (
	"context"
	"time"
)

type SchoolRating struct {
	ID int `json:"id"`

	// Data related to the rating
	Rating        int    `json:"rating"`
	Comment       string `json:"comment"`
	IsApproved    bool   `json:"isApproved" db:"is_approved"`
	ApprovalCount int    `json:"approvalCount" db:"approval_count"`
	UpdatedCount  int    `json:"updateCount" db:"updated_count"`

	// Data related to the school
	School   *School `json:"school,omitempty" db:"-"`
	SchoolID int     `json:"schoolId"`

	User   *User `json:"user"`
	UserID int   `json:"userId" db:"user_id"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type SchoolRatingFilter struct {
	RatingID   *int  `json:"ratingId"`
	SchoolID   *int  `json:"schoolId"`
	IsApproved *bool `json:"isApproved"`
	Offset     int   `json:"offset"`
	Limit      int   `json:"limit"`
}

type SchoolRatingService interface {
	CreateSchoolRating(ctx context.Context, schoolRating *SchoolRating) error

	// Gets school ratings with pagination and filtering
	GetSchoolRatings(ctx context.Context, filter SchoolRatingFilter) ([]*SchoolRating, int, error)

	ApproveSchoolRating(ctx context.Context, id int) error
	UpdateSchoolRating(ctx context.Context, id int, upd *SchoolRatingUpdate) (*SchoolRating, error)
	DeleteSchoolRating(ctx context.Context, id int) error
}

type SchoolRatingUpdate struct {
	Rating  *int    `json:"rating"`
	Comment *string `json:"comment"`
}
