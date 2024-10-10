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

	err = createCountry(ctx, tx, country)
	if err != nil {
		return err
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

	country, err := getCountryById(ctx, tx, id)
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

	country, err := updateCountry(ctx, tx, id, upd)
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
			updated_at,
			flag_code
		from country 
		where ` + strings.Join(where, " and ") + `
		order by id	
	` + FormatLimitOffset(filter.Limit, filter.Offset)

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

func getCountryById(ctx context.Context, tx *Tx, id int) (*etp.Country, error) {
	countries, n, err := getCountries(ctx, tx, etp.CountryFilter{
		CountryId: &id,
	})
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, etp.Errorf(etp.ENOTFOUND, "country with id %d not found", id)
	}

	return countries[0], nil
}

func updateCountry(ctx context.Context, tx *Tx, id int, upd *etp.CountryUpdate) (*etp.Country, error) {
	country, err := getCountryById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if v := upd.Name; v != nil {
		country.Name = *v
	}

	if v := upd.Abbreviation; v != nil {
		country.Abbreviation = *v
	}

	if v := upd.AdditionalFields; v != nil {
		country.AdditionalFields = *v
	}

	query := `
		update country
		set
			name = @name,
			abbreviation = @abbreviation,
			additional_fields = @additionalFields,
			updated_at = now()
		where id = @id
	`

	args := pgx.NamedArgs{
		"id":               id,
		"name":             country.Name,
		"abbreviation":     country.Abbreviation,
		"additionalFields": country.AdditionalFields,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		slog.Error("error while updating country", "error", err)
		return nil, err
	}

	// get time in microseconds golang
	country.UpdatedAt = time.Now()

	return country, err
}

func createCountry(ctx context.Context, tx *Tx, country *etp.Country) error {
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

	_, err := tx.Exec(ctx, query, args)
	if err != nil {
		slog.Error("error while creating country", "country", country, "error", err)
		return err
	}

	return nil
}
