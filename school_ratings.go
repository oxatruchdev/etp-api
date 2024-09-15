package etpapi

import "time"

type SchoolRatings struct {
	ID int `json:"id"`

	// Data related to the rating
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`

	// Data related to the school
	School   *School `json:"school"`
	SchoolID int     `json:"schoolId"`

	IsReviewed bool `json:"isReviewed"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
