package db

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/jackc/pgx/v5"
)

type SchoolService struct {
	db *DB
}

func NewSchoolService(db *DB) *SchoolService {
	return &SchoolService{
		db: db,
	}
}

func (ss *SchoolService) CreateSchool(ctx context.Context, school *etp.School) error {
	tx, err := ss.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		insert into school
			(
				name,
				abbreviation,
				country_id,
				metadata
			)
		values
			(
				@name,
				@abbreviation,
				@countryId,
				@metadata
			)
	`

	args := pgx.NamedArgs{
		"name":         school.Name,
		"abbreviation": school.Abbreviation,
		"countryId":    school.CountryID,
		"metadata":     school.Metadata,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		slog.Error("error while creating school", "school", school, "error", err)
		return err
	}

	return tx.Commit(ctx)
}

func (ss *SchoolService) GetSchoolById(ctx context.Context, id int) (*etp.School, error) {
	tx, err := ss.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	school, err := getSchoolById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return school, nil
}

func (ss *SchoolService) GetSchools(ctx context.Context, filter etp.SchoolFilter) ([]*etp.School, int, error) {
	tx, err := ss.db.BeginTx(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	schools, n, err := getSchools(ctx, tx, filter)
	if err != nil {
		return nil, 0, err
	}

	return schools, n, tx.Commit(ctx)
}

func (ss *SchoolService) UpdateSchool(ctx context.Context, id int, upd *etp.SchoolUpdate) (*etp.School, error) {
	tx, err := ss.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	school, err := updateSchool(ctx, tx, id, upd)
	if err != nil {
		return nil, err
	}

	return school, tx.Commit(ctx)
}

func (ss *SchoolService) DeleteSchool(ctx context.Context, id int) error {
	tx, err := ss.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = deleteSchool(ctx, tx, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func deleteSchool(ctx context.Context, tx *Tx, id int) error {
	_, err := tx.Exec(ctx, "delete from school where id = @id", pgx.NamedArgs{"id": id})
	if err != nil {
		return err
	}
	return nil
}

func updateSchool(ctx context.Context, tx *Tx, id int, upd *etp.SchoolUpdate) (*etp.School, error) {
	school, err := getSchoolById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if v := upd.Name; v != nil {
		school.Name = *v
	}

	if v := upd.Abbreviation; v != nil {
		school.Abbreviation = *v
	}

	if v := upd.Metadata; v != nil {
		school.Metadata = *v
	}

	query := `
		update 
			school
		set
			name = @name,
			abbreviation = @abbreviation,
			metadata = @metadata,
			updated_at = now()
		where
			id = @id
	`

	args := pgx.NamedArgs{
		"id":           id,
		"name":         school.Name,
		"abbreviation": school.Abbreviation,
		"metadata":     school.Metadata,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		slog.Error("error while updating school", "school", school, "error", err)
		return nil, err
	}

	return school, nil
}

func getSchoolById(ctx context.Context, tx *Tx, id int) (*etp.School, error) {
	schools, n, err := getSchools(ctx, tx, etp.SchoolFilter{
		SchoolId: &id,
	})
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, &etp.Error{Code: etp.ENOTFOUND, Message: "school not found"}
	}

	return schools[0], nil
}

func getSchools(ctx context.Context, tx *Tx, filter etp.SchoolFilter) ([]*etp.School, int, error) {
	where, args := []string{"1 = 1"}, pgx.NamedArgs{}

	if filter.SchoolName != nil {
		where = append(where, "unaccent(name) ilike @name or unaccent(abbreviation) ilike @name")
		args["name"] = "%" + *filter.SchoolName + "%"
	}

	if filter.CountryID != nil {
		where = append(where, "country_id = @countryId")
		args["countryId"] = *filter.CountryID
	}

	if filter.SchoolId != nil {
		where = append(where, "id = @id")
		args["id"] = *filter.SchoolId
	}

	query := `
		select 
			count(*)
		from
			school
		where ` + strings.Join(where, " and ")

	var n int
	if err := tx.QueryRow(ctx, query, args).Scan(&n); err != nil {
		return nil, 0, err
	}

	if n == 0 {
		return nil, 0, nil
	}

	query = `
		select
			id,
			name,
			abbreviation,
			metadata,
			country_id,
			created_at,
			updated_at
		from
			school
		where ` + strings.Join(where, " and ") + `
		order by id
		` + FormatLimitOffset(filter.Offset, filter.Limit)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}

	schools, err := pgx.CollectRows(rows, pgx.RowToStructByName[etp.School])
	if err != nil {
		return nil, 0, err
	}

	schoolPtrs := make([]*etp.School, len(schools))
	for i := range schools {
		schoolPtrs[i] = &schools[i]
	}

	return schoolPtrs, n, nil
}
