package db

import (
	"context"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/jackc/pgx/v5"
)

type TagService struct {
	db *DB
}

func NewTagService(db *DB) *TagService {
	return &TagService{
		db: db,
	}
}

func (s *TagService) GetTags(ctx context.Context) ([]*etp.Tag, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	tags, err := getTags(ctx, tx)
	if err != nil {
		return nil, err
	}

	return tags, tx.Commit(ctx)
}

func (s *TagService) CreateTag(ctx context.Context, name string) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = createTag(ctx, tx, name)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func createTag(ctx context.Context, tx *Tx, name string) error {
	query := `
		insert into tag
			(
				name
			)
		values
			(
				@name
			)
	`

	args := pgx.NamedArgs{
		"name": name,
	}

	_, err := tx.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func getTags(ctx context.Context, tx *Tx) ([]*etp.Tag, error) {
	query := `
		select id, name, created_at, updated_at
		from tag
		order by name
	`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	tags, err := pgx.CollectRows(rows, pgx.RowToStructByName[etp.Tag])
	if err != nil {
		return nil, err
	}

	tagsPtrs := make([]*etp.Tag, len(tags))
	for i := range tags {
		tagsPtrs[i] = &tags[i]
	}

	return tagsPtrs, nil
}
