package etpapi

import "time"

type Department struct {
	ID int `json:"id"`

	// Name and abbreviation are used for display
	Name string `json:"name"`

	// Relations
	School   *School `json:"school"`
	SchoolID int     `json:"schoolId"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
