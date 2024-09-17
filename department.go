package etpapi

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
	School     *School      `json:"school"`
	SchoolID   int          `json:"schoolId"`
	Professors *[]Professor `json:"professors"`
	Courses    *[]Course    `json:"courses"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type DepartmentFilter struct {
	ID       *int `json:"id"`
	SchoolID *int `json:"schoolId"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type DepartmentService interface {
	// Gets a department by id
	GetDepartmentById(ctx context.Context, id int) (*Department, error)

	// Gets all departments
	// Offset and Limit are used for pagination
	GetDepartments(ctx context.Context, filter DepartmentFilter) (*[]Department, int, error)

	// Creates a department
	CreateDepartment(ctx context.Context, department *Department) error

	// Updates a department
	UpdateDepartment(ctx context.Context, id int, upd *DepartmentUpdate) (*Department, error)
}

type DepartmentUpdate struct {
	Name *string `json:"name"`
	Code *string `json:"code"`
}
