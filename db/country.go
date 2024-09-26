package db

import (
	"context"
	"log/slog"
	"time"

	etp "github.com/Evalua-Tu-Profe/etp-api"
	"github.com/jackc/pgx/v5"
)

type CountryService struct {
	db *DB
}

func NewCountryService(db *DB) *CountryService {
	return &CountryService{
		db: db,
	}
}

func (cs *CountryService) CreateCountry(ctx context.Context, country *etp.Country) error {
	tx, err := cs.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		insert 
		into country 
			(
				name, 
				abbreviation, 
				additional_fields
			) 
		values 
			(
				@name, 
				@abbreviation, 
				@additionalFields
			)
	`

	args := pgx.NamedArgs{
		"name":             country.Name,
		"abbreviation":     country.Abbreviation,
		"additionalFields": country.AdditionalFields,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		slog.Error("error while creating country", "country", country, "error", err)
	}

	return tx.Commit(ctx)
}

func (cs *CountryService) GetCountries(ctx context.Context) ([]*etp.Country, error) {
	tx, err := cs.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		select 
			id, 
			name, 
			abbreviation, 
			additional_fields, 
			created_at, 
			updated_at 
		from country
	`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	countries, err := pgx.CollectRows(rows, pgx.RowToStructByName[*etp.Country])
	if err != nil {
		return nil, err
	}

	return countries, tx.Commit(ctx)
}

func (cs *CountryService) GetCountryById(ctx context.Context, id int) (*etp.Country, error) {
	tx, err := cs.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	query := `
		select 
			id, 
			name, 
			abbreviation, 
			additional_fields, 
			created_at, 
			updated_at 
		from country 
		where id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	country, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[*etp.Country])
	if err != nil {
		return nil, err
	}

	return country, tx.Commit(ctx)
}

func (cs *CountryService) UpdateCountry(ctx context.Context, id int, upd *etp.CountryUpdate) (*etp.Country, error) {
	tx, err := cs.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	query := `
		update country
		set
			name = @name,
			abbreviation = @abbreviation,
			additional_fields = @additionalFields
			updated_at = @updated_at
		where id = @id
	`

	args := pgx.NamedArgs{
		"id":               id,
		"name":             upd.Name,
		"abbreviation":     upd.Abbreviation,
		"additionalFields": upd.AdditionalFields,
		"updated_at":       time.Now(),
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return nil, err
	}

	country, err := cs.GetCountryById(ctx, id)
	if err != nil {
		return nil, err
	}

	return country, tx.Commit(ctx)
}
