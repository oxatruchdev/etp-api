package etp

import (
	"context"
	"time"
)

type Course struct {
	ID int `json:"id"`

	Name    string `json:"name"`
	Code    string `json:"code"`
	Credits int    `json:"credits"`

	// Relations
	Department   *Department  `json:"department,omitempty" db:"-"`
	DepartmentID int          `json:"departmentId" db:"department_id"`
	School       *School      `json:"school,omitempty" db:"-"`
	SchoolID     int          `json:"schoolId" db:"school_id"`
	Professors   []*Professor `json:"professors,omitempty" db:"-"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CourseFilter struct {
	ID           *int `json:"id"`
	DepartmentId *int `json:"departmentId"`
	SchoolID     *int `json:"schoolID"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type CourseService interface {
	GetCourseById(ctx context.Context, id int) (*Course, error)

	GetCourses(ctx context.Context, filter CourseFilter) ([]*Course, int, error)

	CreateCourse(ctx context.Context, course *Course) error

	UpdateCourse(ctx context.Context, id int, upd *CourseUpdate) (*Course, error)

	DeleteCourse(ctx context.Context, id int) error
}

type CourseUpdate struct {
	Name    *string `json:"name"`
	Code    *string `json:"code"`
	Credits *int    `json:"credits"`
}
