package etpapi

import "time"

type ProfessorRatings struct {
	ID int `json:"id"`

	// Data related to the review itself

	WouldTakeAgain      bool   `json:"wouldTakeAgain"`
	MandatoryAttendance bool   `json:"mandatoryAttendance"`
	Grade               string `json:"grade"`
	TextbookRequired    bool   `json:"textbookRequired"`

	// Relations
	Course   Course `json:"course"`
	CourseId int    `json:"courseId"`

	Professor   Professor `json:"professor"`
	ProfessorId int       `json:"professorId"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
