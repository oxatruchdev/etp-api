package etp

import (
	"context"
	"time"
)

const (
	RoleAdmin       = "admin"
	RoleStudent     = "student"
	RolaModerator   = "moderator"
	RoleSchoolAdmin = "school_admin"
	RoleProfessor   = "professor"
)

type Role struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	DisplayName string `json:"displayName" db:"display_name"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type RoleFilter struct {
	ID   *int
	Name *string
}

type RoleService interface {
	// Create a new role
	CreateRole(ctx context.Context, name string) error

	// Gets a list of roles
	GetRoles(ctx context.Context, filter *RoleFilter) ([]*Role, error)

	// Gets a role by id
	GetRoleById(ctx context.Context, id int) (*Role, error)

	// Gets a role by name
	GetRoleByName(ctx context.Context, name string) (*Role, error)
}
