package db

import (
	"context"
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
		where = append(where, "id = @id")
		args["id"] = *filter.ID
	}

	query := `
		select 
			count(*)
		from professor
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
			id,
			first_name,
			last_name,
			school_id,
			created_at,
			updated_at,
			full_name
		from
			professor
		where ` + strings.Join(where, " and ") + `
		` + FormatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}

	professors, err := pgx.CollectRows(rows, pgx.RowToStructByName[etp.Professor])
	if err != nil {
		return nil, 0, err
	}

	professorPtrs := make([]*etp.Professor, len(professors))
	for i := range professors {
		professorPtrs[i] = &professors[i]
	}

	return professorPtrs, counter, nil
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
