package etp

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID       int         `json:"id" db:"id"`
	Email    string      `json:"email" db:"email" form:"email"`
	Password string      `json:"password" db:"password" form:"password"`
	Name     pgtype.Text `json:"name" db:"name" form:"name"`

	// Relations
	School   *School `json:"school" db:"-"`
	SchoolID *int    `json:"school_id" db:"school_id" form:"school"`

	// Relations
	Role   *Role `json:"role" db:"-"`
	RoleID *int  `json:"role_id" db:"role_id"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type UserService interface {
	// Gets a user by id
	GetUserById(ctx context.Context, id int) (*User, error)

	// Gets a user by email
	GetUserByEmail(ctx context.Context, email string) (*User, error)

	// Registers a user
	RegisterUser(ctx context.Context, user *User) error
}

type UserUpdate struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}
