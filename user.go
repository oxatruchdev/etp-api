package etp

import (
	"context"
	"time"
)

type User struct {
	ID       int     `json:"id" db:"id"`
	Email    string  `json:"email" db:"email"`
	Password string  `json:"password" db:"password"`
	Name     *string `json:"name" db:"name"`

	// Relations
	School   *School `json:"school" db:"-"`
	SchoolID *int    `json:"school_id" db:"school_id"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type UserService interface {
	// Gets a user by id
	GetUserById(ctx context.Context, id int) (*User, error)

	// Registers a user
	RegisterUser(ctx context.Context, user *User) error

	// Updates a user
	UpdateUser(ctx context.Context, id int, upd *UserUpdate) (*User, error)
}

type UserUpdate struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}
