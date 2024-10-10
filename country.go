package etp

import (
	"context"
	"time"
)

// TODO: Add rest of the fields for the school rating
type Country struct {
	ID int `param:"id" query:"id" json:"id"`

	// Name and abbreviation are used for display
	Name             string `json:"name"`
	Abbreviation     string `json:"abbreviation"`
	AdditionalFields any    `json:"additionalFields" db:"additional_fields"`
	FlagCode         string `json:"flag" db:"flag_code"`

	// Relations
	Schools []*School `json:"schools omitempty" db:"-"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CountryFilter struct {
	CountryId *int `param:"id" json:"id"`

	Offset int `json:"offset" query:"offset"`
	Limit  int `json:"limit" query:"limit"`
}

type CountryService interface {
	GetCountryById(ctx context.Context, id int) (*Country, error)

	GetCountries(ctx context.Context, filter CountryFilter) ([]*Country, int, error)

	CreateCountry(ctx context.Context, country *Country) error

	UpdateCountry(ctx context.Context, id int, upd *CountryUpdate) (*Country, error)
}

type CountryUpdate struct {
	Name             *string `json:"name"`
	Abbreviation     *string `json:"abbreviation"`
	AdditionalFields *any    `json:"additionalFields"`
}
