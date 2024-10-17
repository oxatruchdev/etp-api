package etp

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Tag struct {
	ID   int
	Name string

	// Relations
	ProfessorRatings []*ProfessorRating
	SchoolRatings    []*SchoolRating

	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

type TagService interface {
	GetTags(ctx context.Context) ([]*Tag, error)

	CreateTag(ctx context.Context, name string) error
}
