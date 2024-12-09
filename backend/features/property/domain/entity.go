package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSalePriceInvalid      = errors.New("Sale price invalid")
	ErrRentAmountInvalid     = errors.New("Rent Amount Invalid")
	ErrNotValidCurrency      = errors.New("Not valid currency")
	ErrEmptyNeighborhood     = errors.New("neighborhood cannot be empty")
	ErrEmptyCity             = errors.New("city cannot be empty")
	ErrEmptyDescription      = errors.New("description cannot be empty")
	ErrNegativeBedrooms      = errors.New("number of bedrooms cannot be negative")
	ErrNegativeBathrooms     = errors.New("number of bathrooms cannot be negative")
	ErrInvalidArea           = errors.New("area must be greater than 0")
	ErrInvalidLatitude       = errors.New("latitude must be between -90 and 90")
	ErrInvalidLongitude      = errors.New("longitude must be between -180 and 180")
	ErrInvalidSalePrice      = errors.New("sale price must be 0 if the property is not for sale")
	ErrNotValidTypeProperty  = errors.New("type property not valid, must be 'apartment' or 'house'")
	ErrNotForRent            = errors.New("not for rent")
	ErrInvalidSearchCriteria = errors.New("at least one of 'is_for_rent' or 'is_for_sale' must be true")
)

// later convert the currency  use the api: https://docs.awesomeapi.com.br/api-de-moedas

// get geolocation with https://github.com/kellydunn/golang-geo?tab=readme-ov-fil

type SearchFilter struct {
	PropertyType string
	Neighborhood string
	City         string
	ForRent      bool
	ForSale      bool
	RentAmount   int
	SalePrice    int
	Size         float64
}

// price allways in cents
type Price struct {
	Amount   int
	Currency string
}

type Property struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	PropertyType string     `json:"property_type"`
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
	Code         string     `json:"code"`
	IsForRent    bool       `json:"is_for_rent"`
}

var propertyTypes = map[string]bool{
	"apartment": true,
	"house":     true,
}

// price allways in cents

func NewProperty(neighborhood, city, description string, rent Price, bedrooms, bathrooms int, area, latitude, longitude float64, isForSale bool, salePrice int, propertyType string, isForRent bool) (*Property, error) {
	if !isForRent && !isForSale {
		return nil, ErrInvalidSearchCriteria
	}
	if rent.Amount < 0 && !isForRent {
		return nil, ErrNotForRent
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
	if rent.Amount < 0 && isForSale {
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
	if _, ok := propertyTypes[propertyType]; !ok {
		fmt.Println(propertyType)
		return nil, ErrNotValidTypeProperty
	}
	propertyCode, err := createPropertyCode()
	if err != nil {
		return nil, err
	}

	propertyName := fmt.Sprintf("%s-%s-%s-%d-%d", propertyType, city, neighborhood, bedrooms, bathrooms)

	id, _ := uuid.NewRandom()

	p := &Property{
		ID:           id,
		Name:         propertyName,
		Neighborhood: neighborhood,
		PropertyType: propertyType,
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
		Code:         propertyCode,
		IsForRent:    isForRent,
	}
	return p, nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func createPropertyCode() (string, error) {
	maxLength := 10
	b := make([]byte, maxLength)
	n, err := io.ReadAtLeast(rand.Reader, b, maxLength)
	if n != maxLength {
		return "", err
	}

	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b), nil
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'A', 'B', 'C', 'D', 'E', 'F'}
