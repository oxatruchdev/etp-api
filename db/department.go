package db

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/jackc/pgx/v5"
)

type DepartmentService struct {
	db *DB
}

func NewDepartmentService(db *DB) *DepartmentService {
	return &DepartmentService{
		db: db,
	}
}

func (ds *DepartmentService) CreateDepartment(ctx context.Context, department *etp.Department) error {
	tx, err := ds.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = createDepartment(ctx, tx, department)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (ds *DepartmentService) GetDepartmentById(ctx context.Context, id int) (*etp.Department, error) {
	tx, err := ds.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	departments, n, err := getDepartments(ctx, tx, &etp.DepartmentFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, &etp.Error{Code: etp.ENOTFOUND, Message: "department not found"}
	}

	return departments[0], tx.Commit(ctx)
}

func (ds *DepartmentService) GetDepartments(ctx context.Context, filter etp.DepartmentFilter) ([]*etp.Department, int, error) {
	slog.Debug("GetDepartments", "filter", filter)
	tx, err := ds.db.BeginTx(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	departments, n, err := getDepartments(ctx, tx, &filter)
	if err != nil {
		return nil, 0, err
	}

	if n == 0 {
		return []*etp.Department{}, 0, nil
	}

	return departments, n, tx.Commit(ctx)
}

func (ds *DepartmentService) UpdateDepartment(ctx context.Context, id int, upd *etp.DepartmentUpdate) (*etp.Department, error) {
	tx, err := ds.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	department, err := updateDepartment(ctx, tx, id, upd)
	if err != nil {
		return nil, err
	}

	return department, tx.Commit(ctx)
}

func updateDepartment(ctx context.Context, tx *Tx, id int, upd *etp.DepartmentUpdate) (*etp.Department, error) {
	department, err := getDepartmenById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if v := upd.Name; v != nil {
		department.Name = *v
	}

	if v := upd.Code; v != nil {
		department.Code = *v
	}

	query := `
		update
			department
		set
			name = @name,
			code = @code,
			updated_at = now()
		where
			id = @id
	`

	args := pgx.NamedArgs{
		"id":   id,
		"name": department.Name,
		"code": department.Code,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		slog.Error("Error updating department", "err", err)
		return nil, err
	}

	return department, nil
}

func getDepartmenById(ctx context.Context, tx *Tx, id int) (*etp.Department, error) {
	departments, _, err := getDepartments(ctx, tx, &etp.DepartmentFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	if len(departments) == 0 {
		return nil, &etp.Error{Code: etp.ENOTFOUND, Message: "department not found"}
	}

	return departments[0], nil
}

func getDepartments(ctx context.Context, tx *Tx, filter *etp.DepartmentFilter) ([]*etp.Department, int, error) {
	where, args := []string{"1 = 1"}, pgx.NamedArgs{}

	if filter.ID != nil {
		where = append(where, "id = @id")
		args["id"] = *filter.ID
	}

	if filter.SchoolID != nil {
		where = append(where, "school_id = @schoolId")
		args["schoolId"] = *filter.SchoolID
	}

	query := `
		select 
			count(*)
		from department
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
			name,
			code,
			school_id,
			created_at,
			updated_at
		from
			department
		where ` + strings.Join(where, " and ") + `
		order by name asc
		` + FormatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}

	departments, err := pgx.CollectRows(rows, pgx.RowToStructByName[etp.Department])
	if err != nil {
		slog.Error("error while getting departments", "error", err)
		return nil, 0, err
	}

	departmentsPtr := make([]*etp.Department, len(departments))
	for i := range departments {
		departmentsPtr[i] = &departments[i]
	}

	return departmentsPtr, counter, nil
}

func createDepartment(ctx context.Context, tx *Tx, department *etp.Department) error {
	query := `
		insert into department
			(
				name,
				code,
				school_id
			)
		values
			(
				@name,
				@code,
				@schoolId
			)
	`

	args := pgx.NamedArgs{
		"name":     department.Name,
		"schoolId": department.SchoolID,
		"code":     department.Code,
	}

	_, err := tx.Exec(ctx, query, args)
	if err != nil {
		slog.Error("error while creating department", "department", department, "error", err)
		return err
	}

	return nil
}
