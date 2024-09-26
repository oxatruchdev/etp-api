package etp

import (
	"context"
	"time"
)

// TODO: Add rest of the fields for the school rating
type Country struct {
	ID int `json:"id"`

	// Name and abbreviation are used for display
	Name             string `json:"name"`
	Abbreviation     string `json:"abbreviation"`
	AdditionalFields any    `json:"additionalFields" db:"additional_fields"`

	// Relations
	Schools []*School `json:"schools"`

	// CreatedAt and UpdatedAt are used for tracking
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CountryFilter struct {
	CountryId *int `json:"id"`

	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type CountryService interface {
	GetCountryById(ctx context.Context, id int) (*Country, error)

	GetCountries(ctx context.Context, filter CountryFilter) ([]*Country, int, error)

	CreateCountry(ctx context.Context, country *Country) error

	UpdateCountry(ctx context.Context, id int, upd *CountryUpdate) (*Country, error)
}

type CountryUpdate struct {
	Name             *string `json:"name"`
	Abbreviation     string  `json:"abbreviation"`
	AdditionalFields any     `json:"additionalFields"`
}
