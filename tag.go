package etp

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	// Relations
	ProfessorRatings []*ProfessorRating `json:"professorRatings" db:"-"`
	SchoolRatings    []*SchoolRating    `json:"schoolRatings" db:"-"`

	CreatedAt pgtype.Timestamptz `json:"createdAt" db:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt" db:"updated_at"`
}

type TagWithCount struct {
	Tag
	Count int `json:"count" db:"usage_count"`
}

type TagService interface {
	GetTags(ctx context.Context) ([]*Tag, error)

	CreateTag(ctx context.Context, name string) error
}
