package db

import (
	"context"
	"log/slog"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *DB
}

func NewUserService(db *DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, user *etp.User) error {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = createUser(ctx, tx, user)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *UserService) GetUserById(ctx context.Context, id int) (*etp.User, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	user, err := getUserById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return user, tx.Commit(ctx)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*etp.User, error) {
	tx, err := s.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	user, err := getUserByEmail(ctx, tx, email)
	if err != nil {
		return nil, err
	}

	return user, tx.Commit(ctx)
}

func getUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*etp.User, error) {
	query := `
		SELECT *
		FROM "user"
		WHERE email = @email
	`

	rows, err := tx.Query(ctx, query, pgx.NamedArgs{"email": email})
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[etp.User])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func createUser(ctx context.Context, tx pgx.Tx, user *etp.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	query := `
		INSERT INTO "user"
		(email, password, role_id)
		VALUES
		(@email, @password, @role_id)
	`

	args := pgx.NamedArgs{
		"email":    user.Email,
		"password": user.Password,
		"role_id":  user.RoleID,
	}

	slog.Info("Inserting user", "query", query, "args", args)

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func getUserById(ctx context.Context, tx pgx.Tx, id int) (*etp.User, error) {
	query := `
		SELECT *
		FROM "user"
		WHERE id = @id
	`

	rows, err := tx.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[etp.User])
	if err != nil {
		return nil, err
	}

	return &user, nil
}
