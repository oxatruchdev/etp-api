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

func (s *SchoolService) GetSchoolProfessorsCount(ctx context.Context, id int) (int, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return 0, err
	}

	defer tx.Rollback(ctx)
	count, err := getSchoolProfessorsCount(ctx, tx, id)
	if err != nil {
		return 0, err
	}

	return count, tx.Commit(ctx)
}

func getSchoolProfessorsCount(ctx context.Context, tx *Tx, id int) (int, error) {
	var count int

	slog.Info("Getting school professors count for school: ", "id", id)
	err := tx.QueryRow(ctx, "select count(*) from professor where school_id = @id", pgx.NamedArgs{"id": id}).Scan(&count)
	if err != nil {
		slog.Error("error while getting school professors count", "error", err, "school id", id)
		return 0, err
	}

	return count, nil
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
		where = append(where, "unaccent(s.name) ilike @name or unaccent(s.abbreviation) ilike @name")
		args["name"] = "%" + *filter.SchoolName + "%"
	}

	if filter.CountryID != nil {
		where = append(where, "country_id = @countryId")
		args["countryId"] = *filter.CountryID
	}

	if filter.SchoolId != nil {
		where = append(where, "s.id = @id")
		args["id"] = *filter.SchoolId
	}

	query := `
		select 
			count(*)
		from
			school s
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
			s.id,
			s.name,
			s.abbreviation,
			metadata,
			country_id,
			s.created_at,
			s.updated_at,
			c.id as country_id,
			c.name as country_name,
			c.abbreviation as country_abbreviation,
			c.flag_code
		from
			school s
		left join country c on c.id = country_id
		where ` + strings.Join(where, " and ") + `
		order by random()
		` + FormatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}

	var schools []etp.School
	for rows.Next() {
		var school etp.School
		var country etp.Country
		if err := rows.Scan(
			&school.ID,
			&school.Name,
			&school.Abbreviation,
			&school.Metadata,
			&school.CountryID,
			&school.CreatedAt,
			&school.UpdatedAt,
			&country.ID,
			&country.Name,
			&country.Abbreviation,
			&country.FlagCode,
		); err != nil {
			return nil, 0, err
		}

		school.Country = &country
		schools = append(schools, school)
	}

	schoolPtrs := make([]*etp.School, len(schools))
	for i := range schools {
		schoolPtrs[i] = &schools[i]
	}

	return schoolPtrs, n, nil
}
