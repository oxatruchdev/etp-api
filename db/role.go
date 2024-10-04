package db

import (
	"context"
	"strings"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/jackc/pgx/v5"
)

type RoleService struct {
	db *DB
}

func NewRoleService(db *DB) *RoleService {
	return &RoleService{
		db: db,
	}
}

func (s *RoleService) CreateRole(ctx context.Context, name string) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = createRole(ctx, tx, name)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *RoleService) GetRoles(ctx context.Context, filter *etp.RoleFilter) ([]*etp.Role, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	roles, err := getRoles(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	return roles, tx.Commit(ctx)
}

func (s *RoleService) GetRoleById(ctx context.Context, id int) (*etp.Role, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	role, err := getRoleById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return role, tx.Commit(ctx)
}

func (s *RoleService) GetRoleByName(ctx context.Context, name string) (*etp.Role, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	role, err := getRoleByName(ctx, tx, name)
	if err != nil {
		return nil, err
	}

	return role, tx.Commit(ctx)
}

func createRole(ctx context.Context, tx *Tx, name string) error {
	query := `
		INSERT INTO "role"
		(name)
		VALUES
		(@name)
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

func getRoleById(ctx context.Context, tx *Tx, id int) (*etp.Role, error) {
	roles, err := getRoles(ctx, tx, &etp.RoleFilter{
		ID: &id,
	})
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return nil, &etp.Error{Code: etp.ENOTFOUND, Message: "role not found"}
	}

	return roles[0], nil
}

func getRoleByName(ctx context.Context, tx *Tx, name string) (*etp.Role, error) {
	roles, err := getRoles(ctx, tx, &etp.RoleFilter{
		Name: &name,
	})
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return nil, &etp.Error{Code: etp.ENOTFOUND, Message: "role not found"}
	}

	return roles[0], nil
}

func getRoles(ctx context.Context, tx *Tx, filter *etp.RoleFilter) ([]*etp.Role, error) {
	where, args := []string{"1 = 1"}, pgx.NamedArgs{}

	if v := filter.ID; v != nil {
		where = append(where, "id = @id")
		args["id"] = *v
	}

	if v := filter.Name; v != nil {
		where = append(where, "name = @name")
		args["name"] = *v
	}

	query := `
		SELECT *
		FROM role
		WHERE ` + strings.Join(where, " AND ")

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	roles, err := pgx.CollectRows(rows, pgx.RowToStructByName[etp.Role])
	if err != nil {
		return nil, err
	}

	rolesPtr := make([]*etp.Role, len(roles))

	for i := range roles {
		rolesPtr[i] = &roles[i]
	}

	return rolesPtr, nil
}
