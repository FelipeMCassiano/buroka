package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSalePriceInvalid  = errors.New("Sale price invalid")
	ErrRentAmountInvalid = errors.New("Rent Amount Invalid")
	ErrNotValidCurrency  = errors.New("Not valid currency")
	ErrEmptyName         = errors.New("property name cannot be empty")
	ErrEmptyNeighborhood = errors.New("neighborhood cannot be empty")
	ErrEmptyCity         = errors.New("city cannot be empty")
	ErrEmptyDescription  = errors.New("description cannot be empty")
	ErrNegativeBedrooms  = errors.New("number of bedrooms cannot be negative")
	ErrNegativeBathrooms = errors.New("number of bathrooms cannot be negative")
	ErrInvalidArea       = errors.New("area must be greater than 0")
	ErrInvalidLatitude   = errors.New("latitude must be between -90 and 90")
	ErrInvalidLongitude  = errors.New("longitude must be between -180 and 180")
	ErrInvalidSalePrice  = errors.New("sale price must be 0 if the property is not for sale")
)

// later convert the currency  use the api: https://docs.awesomeapi.com.br/api-de-moedas

// get geolocation with https://github.com/kellydunn/golang-geo?tab=readme-ov-fil

// price allways in cents
type Price struct {
	Amount   int
	Currency string
}

type Property struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Neighborhood string     `json:"neighborhood"`
	City         string     `json:"city"`
	Rent         Price      `json:"rent"`
	Bedrooms     int        `json:"bedrooms"`
	Bathrooms    int        `json:"bathrooms"`
	Area         float64    `json:"area"`
	Description  string     `json:"description"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	CreatedAt    *time.Time `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	IsForSale    bool       `json:"is_for_sale"`
	SalePrice    int        `json:"sale_price"`
}

// price allways in cents

func NewProperty(name, neighborhood, city, description string, rent Price, bedrooms, bathrooms int, area, latitude, longitude float64, isForSale bool, salePrice int) (*Property, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	if neighborhood == "" {
		return nil, ErrEmptyNeighborhood
	}
	if city == "" {
		return nil, ErrEmptyCity
	}
	if description == "" {
		return nil, ErrEmptyDescription
	}
	if rent.Amount < 0 {
		return nil, ErrRentAmountInvalid
	}
	if rent.Currency == "" {
		return nil, ErrNotValidCurrency
	}
	if bedrooms < 0 {
		return nil, ErrNegativeBedrooms
	}
	if bathrooms < 0 {
		return nil, ErrNegativeBathrooms
	}
	if area <= 0 {
		return nil, ErrInvalidArea
	}
	if latitude < -90 || latitude > 90 {
		return nil, ErrInvalidLatitude
	}
	if longitude < -180 || longitude > 180 {
		return nil, ErrInvalidLongitude
	}
	if isForSale && salePrice < 0 {
		return nil, ErrSalePriceInvalid
	}
	if !isForSale && salePrice != 0 {
		return nil, ErrInvalidSalePrice
	}
	id, _ := uuid.NewRandom()

	p := &Property{
		ID:           id,
		Name:         sanitizeName(name),
		Neighborhood: neighborhood,
		City:         city,
		Description:  description,
		Rent:         rent,
		Bedrooms:     bedrooms,
		Bathrooms:    bathrooms,
		Area:         area,
		Latitude:     latitude,
		Longitude:    longitude,
		IsForSale:    isForSale,
		SalePrice:    salePrice,
		CreatedAt:    timePtr(time.Now()),
	}
	return p, nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func sanitizeName(name string) string {
	return strings.ReplaceAll(name, " ", "-")
}
