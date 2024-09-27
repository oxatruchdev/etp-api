package etp

import (
	"context"
	"time"
)

type Department struct {
	ID int `json:"id"`

	// Department properties
	Name string `json:"name"`
	Code string `json:"code"`

	// Relations
	School     *School      `json:"school,omitempty" db:"-"`
	SchoolID   int          `json:"schoolId" db:"school_id"`
	Professors []*Professor `json:"professors,omitempty" db:"-"`
	Courses    []*Course    `json:"courses,omitempty" db:"-"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type DepartmentFilter struct {
	ID       *int `json:"id" query:"id"`
	SchoolID *int `json:"schoolId" query:"schoolId"`

	Offset int `json:"offset" query:"offset"`
	Limit  int `json:"limit" query:"limit"`
}

type DepartmentService interface {
	// Gets a department by id
	GetDepartmentById(ctx context.Context, id int) (*Department, error)

	// Gets all departments
	// Offset and Limit are used for pagination
	GetDepartments(ctx context.Context, filter DepartmentFilter) ([]*Department, int, error)

	// Creates a department
	CreateDepartment(ctx context.Context, department *Department) error

	// Updates a department
	UpdateDepartment(ctx context.Context, id int, upd *DepartmentUpdate) (*Department, error)
}

type DepartmentUpdate struct {
	Name *string `json:"name"`
	Code *string `json:"code"`
}
