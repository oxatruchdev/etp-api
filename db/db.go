package db

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type DB struct {
	db     *pgxpool.Pool
	ctx    context.Context
	cancel func()

	Now func() time.Time
	DSN string
}

func NewDB(dsn string) *DB {
	db := &DB{
		DSN: dsn,
		Now: time.Now,
	}

	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Ping() error {
	// pgxpool does not have a direct Ping method, we use Exec for a simple query to check connection health.
	conn, err := db.db.Acquire(db.ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	err = conn.QueryRow(db.ctx, "SELECT 1").Scan(new(int))
	return err
}

func (db *DB) Open() (err error) {
	if db.DSN == "" {
		return fmt.Errorf("dsn is required")
	}

	db.db, err = pgxpool.New(db.ctx, db.DSN)
	if err != nil {
		return err
	}

	if err := db.migrate(); err != nil {
		return err
	}

	return nil
}

// migrate sets up migration tracking and executes pending migration files.
func (db *DB) migrate() error {
	// Ensure the 'migrations' table exists so we don't duplicate migrations.
	if _, err := db.db.Exec(db.ctx, `CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	// Read migration files from our embedded file system.
	names, err := fs.Glob(migrationFS, "migration/*.sql")
	if err != nil {
		return err
	}

	slices.Sort(names)

	// Loop over all migration files and execute them in order.
	for _, name := range names {
		if err := db.migrateFile(name); err != nil {
			return fmt.Errorf("migration error: name=%q err=%w", name, err)
		}
	}
	return nil
}

// migrateFile runs a single migration file within a transaction.
func (db *DB) migrateFile(name string) error {
	tx, err := db.db.Begin(db.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(db.ctx)

	// Ensure migration has not already been run.
	var n int
	if err := tx.QueryRow(db.ctx, `SELECT COUNT(*) FROM migrations WHERE name = $1`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil // already run migration, skip
	}

	// Read and execute migration file.
	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(db.ctx, string(buf)); err != nil {
		return err
	}

	// Insert record into migrations to prevent re-running migration.
	if _, err := tx.Exec(db.ctx, `INSERT INTO migrations (name) VALUES (@name)`, pgx.NamedArgs{"name": name}); err != nil {
		return err
	}

	return tx.Commit(db.ctx)
}

// Close closes the database connection.
func (db *DB) Close() error {
	db.cancel()   // Cancel the context
	db.db.Close() // Close the pgxpool connection
	return nil
}

// BeginTx starts a transaction and returns a wrapper Tx type.
func (db *DB) BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := db.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now().UTC().Truncate(time.Second),
	}, nil
}

// Tx is a wrapper around a pgx.Tx transaction.
type Tx struct {
	pgx.Tx
	db  *DB
	now time.Time
}

// FormatLimitOffset returns a SQL string for a given limit & offset.
func FormatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)
	} else if offset > 0 {
		return fmt.Sprintf(`OFFSET %d`, offset)
	}
	return ""
}
