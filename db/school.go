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
				country_id
			)
		values
			(
				@name,
				@abbreviation,
				@countryId
			)
	`

	args := pgx.NamedArgs{
		"name":         school.Name,
		"abbreviation": school.Abbreviation,
		"countryId":    school.CountryID,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		slog.Error("error while creating school", school, err)
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

	query := `
		update 
			school
		set
			name = @name,
			abbreviation = @abbreviation
			updated_at = now()
		where
			id = @id
	`

	args := pgx.NamedArgs{
		"id":           id,
		"name":         school.Name,
		"abbreviation": school.Abbreviation,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
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

	if filter.CountryID != nil {
		where = append(where, "country_id = @countryId")
		args["countryId"] = *filter.CountryID
	}

	if filter.SchoolId != nil {
		where = append(where, "id = @id")
		args["id"] = *filter.SchoolId
	}

	var n int
	query := `
		select
			id,
			name,
			abbreviation,
			country_id
			created_at,
			updated_at,
			count(*) over()
		from
			school
		where ` + strings.Join(where, " and ") + `
		` + FormatLimitOffset(filter.Offset, filter.Limit)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var schools []*etp.School
	for rows.Next() {
		var school *etp.School
		if err := rows.Scan(
			&school.ID,
			&school.Name,
			&school.Abbreviation,
			&school.CountryID,
			&school.CreatedAt,
			&school.UpdatedAt,
			&n,
		); err != nil {
			return nil, 0, err
		}

		schools = append(schools, school)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return schools, n, nil
}
