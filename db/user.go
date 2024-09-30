package db

import (
	"context"

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

func createUser(ctx context.Context, tx pgx.Tx, user *etp.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	query := `
		INSERT INTO "user"
		(email, password)
		VALUES
		(@email, @password)
	`

	args := pgx.NamedArgs{
		"email":    user.Email,
		"password": user.Password,
	}

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

	var user etp.User
	err := tx.QueryRow(ctx, query, pgx.NamedArgs{"id": id}).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
