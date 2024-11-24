package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/jackc/pgx/v5"
)

type ProfessorService struct {
	db *DB
}

func NewProfessorService(db *DB) *ProfessorService {
	return &ProfessorService{
		db: db,
	}
}

func (ps *ProfessorService) CreateProfessor(ctx context.Context, professor *etp.Professor) error {
	tx, err := ps.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = createProfessor(ctx, tx, professor)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (ps *ProfessorService) GetProfessorById(ctx context.Context, id int) (*etp.Professor, error) {
	tx, err := ps.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	professors, n, err := getProfessors(ctx, tx, &etp.ProfessorFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, &etp.Error{Code: etp.ENOTFOUND, Message: "professor not found"}
	}

	return professors[0], tx.Commit(ctx)
}

func (ps *ProfessorService) GetProfessors(ctx context.Context, filter etp.ProfessorFilter) ([]*etp.Professor, int, error) {
	tx, err := ps.db.BeginTx(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	professors, n, err := getProfessors(ctx, tx, &filter)
	if err != nil {
		return nil, 0, err
	}

	if n == 0 {
		return nil, 0, nil
	}

	return professors, n, tx.Commit(ctx)
}

func (ps *ProfessorService) UpdateProfessor(ctx context.Context, id int, upd *etp.ProfessorUpdate) (*etp.Professor, error) {
	tx, err := ps.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	professor, err := updateProfessor(ctx, tx, id, upd)
	if err != nil {
		return nil, err
	}

	return professor, tx.Commit(ctx)
}

func (ps *ProfessorService) DeleteProfessor(ctx context.Context, id int) error {
	tx, err := ps.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := deleteProfessor(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *ProfessorService) GetProfessorCourses(ctx context.Context, id int) ([]*etp.Course, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return []*etp.Course{}, err
	}
	defer tx.Rollback(ctx)

	courses, err := getProfessorCourses(ctx, tx, id)
	if err != nil {
		return []*etp.Course{}, err
	}

	return courses, tx.Commit(ctx)
}

func (s *ProfessorService) GetProfessorTags(ctx context.Context, id int) ([]*etp.TagWithCount, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return []*etp.TagWithCount{}, err
	}
	defer tx.Rollback(ctx)

	tags, err := getProfessorPopularTags(ctx, tx, id)
	if err != nil {
		return []*etp.TagWithCount{}, err
	}

	return tags, tx.Commit(ctx)
}

func (s *ProfessorService) GetLatestProfessors(ctx context.Context, filter etp.ProfessorFilter) ([]*etp.Professor, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return []*etp.Professor{}, err
	}
	defer tx.Rollback(ctx)

	professors, err := getLatestProfessorsWithRatings(ctx, tx, filter)
	if err != nil {
		return []*etp.Professor{}, err
	}

	return professors, tx.Commit(ctx)
}

func getLatestProfessorsWithRatings(ctx context.Context, tx *Tx, filter etp.ProfessorFilter) ([]*etp.Professor, error) {
	return nil, nil
}

func getProfessorPopularTags(ctx context.Context, tx *Tx, id int) ([]*etp.TagWithCount, error) {
	query := `
		WITH tag_counts AS (
			SELECT 
					t.id,
					t.name,
					pr.professor_id,
					COUNT(*) as usage_count
			FROM professor_rating_tag prt
			JOIN tag t ON t.id = prt.tag_id
			JOIN professor_rating pr ON pr.id = prt.professor_rating_id
			WHERE pr.professor_id = @professorId
			AND pr.is_approved = true
			GROUP BY t.id, t.name, pr.professor_id
			ORDER BY usage_count DESC
			LIMIT 3
		)
		SELECT 
				COALESCE(jsonb_agg(
						jsonb_build_object(
								'id', id,
								'name', name,
								'usage_count', usage_count
						)
						ORDER BY usage_count DESC
						
				), '[]'::jsonb) as popular_tags
		FROM tag_counts
	`

	var tagsJSON []byte
	err := tx.QueryRow(ctx, query, pgx.NamedArgs{
		"professorId": id,
	}).Scan(&tagsJSON)
	if err != nil {
		if err == pgx.ErrNoRows {
			return []*etp.TagWithCount{}, nil // return empty slice if no results
		}
		return nil, err
	}

	var popularTags []*etp.TagWithCount
	if err := json.Unmarshal(tagsJSON, &popularTags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal popular tags: %w", err)
	}

	return popularTags, nil
}

func getProfessorCourses(ctx context.Context, tx *Tx, id int) ([]*etp.Course, error) {
	query := `
		select 
			c.id,
			name,
			code
		from course c
		left join professor_course pc on pc.course_id = c.id
		where pc.professor_id = @professorId;	
	`

	args := pgx.NamedArgs{
		"professorId": id,
	}

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return []*etp.Course{}, err
	}
	defer rows.Close()

	var course []*etp.Course
	for rows.Next() {
		var department etp.Course
		err := rows.Scan(
			&department.ID,
			&department.Name,
			&department.Code,
		)
		if err != nil {
			return []*etp.Course{}, err
		}
		course = append(course, &department)
	}

	return course, nil
}

func updateProfessor(ctx context.Context, tx *Tx, id int, upd *etp.ProfessorUpdate) (*etp.Professor, error) {
	professor, err := getProfessorById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if v := upd.FirstName; v != nil {
		professor.FirstName = *v
	}

	if v := upd.LastName; v != nil {
		professor.LastName = *v
	}

	query := `
		update
			professor
		set
			first_name = @firstName,
			last_name = @lastName,
			updated_at = now()
		where
			id = @id
	`

	args := pgx.NamedArgs{
		"id":        id,
		"firstName": professor.FirstName,
		"lastName":  professor.LastName,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return professor, nil
}

func getProfessorById(ctx context.Context, tx *Tx, id int) (*etp.Professor, error) {
	professors, _, err := getProfessors(ctx, tx, &etp.ProfessorFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	if len(professors) == 0 {
		return nil, &etp.Error{Code: etp.ENOTFOUND, Message: "professor not found"}
	}

	return professors[0], nil
}

func getProfessors(ctx context.Context, tx *Tx, filter *etp.ProfessorFilter) ([]*etp.Professor, int, error) {
	where, args := []string{"1 = 1"}, pgx.NamedArgs{}

	if filter.Name != nil {
		where = append(where, "unaccent(first_name) ilike @name or unaccent(last_name) ilike @name or unaccent(full_name) ilike @name")
		args["name"] = "%" + *filter.Name + "%"
	}

	if filter.ID != nil {
		where = append(where, "p.id = @id")
		args["id"] = *filter.ID
	}

	query := `
		select 
			count(*)
		from professor p
		where ` + strings.Join(where, " and ")

	var counter int
	err := tx.QueryRow(ctx, query, args).Scan(&counter)
	if err != nil {
		return nil, 0, err
	}

	if counter == 0 {
		return nil, 0, nil
	}

	query = `
		select
			p.id,
			first_name,
			last_name,
			p.school_id,
			p.created_at,
			p.updated_at,
			full_name,
			d.id as department_id,
			d.name,
			d.code
		from
			professor p
		join department d on d.id = p.department_id
		where ` + strings.Join(where, " and ") + `
		` + FormatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	professors := make([]*etp.Professor, 0)

	for rows.Next() {
		professor := etp.Professor{}
		department := etp.Department{}
		err = rows.Scan(
			&professor.ID,
			&professor.FirstName,
			&professor.LastName,
			&professor.SchoolId,
			&professor.CreatedAt,
			&professor.UpdatedAt,
			&professor.FullName,
			&department.ID,
			&department.Name,
			&department.Code,
		)
		if err != nil {
			return nil, 0, err
		}
		professor.Department = department
		professors = append(professors, &professor)
	}

	if rows.Err() != nil {
		return nil, 0, rows.Err()
	}

	return professors, counter, nil
}

func createProfessor(ctx context.Context, tx *Tx, professor *etp.Professor) error {
	query := `
		insert into professor
			(
				first_name,
				last_name,
				school_id
			)
		values
			(
				@firstName,
				@lastName,
				@schoolId
			)
	`

	args := pgx.NamedArgs{
		"firstName": professor.FirstName,
		"lastName":  professor.LastName,
		"schoolId":  professor.SchoolId,
	}

	_, err := tx.Exec(ctx, query, args)
	if err != nil {
		slog.Error("error while creating professor", "professor", professor, "error", err)
		return err
	}

	return nil
}

func deleteProfessor(ctx context.Context, tx *Tx, id int) error {
	query := `
		delete from professor
		where id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := tx.Exec(ctx, query, args)
	if err != nil {
		slog.Error("error while deleting professor", "id", id, "error", err)
		return err
	}

	return nil
}
