package etpapi

import "time"

type Professor struct {
	ID int `json:"id"`

	// Personal Information
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	School   *School `json:"school"`
	SchoolId int     `json:"schoolId"`

	// Relations
	Ratings *[]ProfessorRatings `json:"ratings"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
