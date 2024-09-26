package db

import (
	"context"
	"log/slog"
	"strings"
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

func (cs *CountryService) GetCountries(ctx context.Context, filter etp.CountryFilter) ([]*etp.Country, int, error) {
	tx, err := cs.db.BeginTx(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	countries, n, err := getCountries(ctx, tx, filter)
	if err != nil {
		return nil, 0, err
	}

	if n == 0 {
		return []*etp.Country{}, 0, nil
	}

	return countries, n, tx.Commit(ctx)
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

func getCountries(ctx context.Context, tx *Tx, filter etp.CountryFilter) ([]*etp.Country, int, error) {
	where, args := []string{"1=1"}, pgx.NamedArgs{}

	if v := filter.CountryId; v != nil {
		where = append(where, "id = @id")
		args["id"] = *v
	}

	countQuery := `
		select 
			count(*) 
		from country
		where ` + strings.Join(where, " and ")

	slog.Info("Countries count query", "query", countQuery, "args", args)

	var n int
	err := tx.QueryRow(ctx, countQuery, args).Scan(&n)
	if err != nil {
		return []*etp.Country{}, 0, err
	}

	query := `
		select 
			id,
			name,
			abbreviation,
			additional_fields,
			created_at,
			updated_at
		from country 
		where ` + strings.Join(where, " and ") + `
	` + FormatLimitOffset(filter.Limit, filter.Offset)

	slog.Info("Countries query", "query", query, "args", args)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}

	countries, err := pgx.CollectRows(rows, pgx.RowToStructByName[etp.Country])
	if err != nil {
		return nil, 0, err
	}

	// Convert the slice of values to a slice of pointers.
	countryPtrs := make([]*etp.Country, len(countries))
	for i := range countries {
		countryPtrs[i] = &countries[i] // take address of each country
	}
	return countryPtrs, n, nil
}
