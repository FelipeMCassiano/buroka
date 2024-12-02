package application

import (
	"fmt"

	"github.com/FelipeMCassiano/buroka/features/property/domain"
	"github.com/FelipeMCassiano/buroka/features/property/infrastructure"
)

type PropertyService struct {
	repo *infrastructure.PropertyRepository
}

type RegisterPropertyRequest struct {
	Name         string       `json:"name"`
	Neighborhood string       `json:"neighborhood"`
	City         string       `json:"city"`
	Rent         domain.Price `json:"rent"`
	Bedrooms     int          `json:"bedrooms"`
	Bathrooms    int          `json:"bathrooms"`
	Area         float64      `json:"area"`
	Description  string       `json:"description"`
	Latitude     float64      `json:"latitude"`
	Longitude    float64      `json:"longitude"`
	IsForSale    bool         `json:"is_for_sale"`
	SalePrice    int          `json:"sale_price"`
}

func NewPropertyService(repo *infrastructure.PropertyRepository) *PropertyService {
	return &PropertyService{repo: repo}
}

func (s *PropertyService) RegisterProperty(propertyRequest *RegisterPropertyRequest) error {
	property, err := domain.NewProperty(propertyRequest.Name, propertyRequest.Neighborhood, propertyRequest.City, propertyRequest.Description, propertyRequest.Rent, propertyRequest.Bedrooms, propertyRequest.Bathrooms, propertyRequest.Area, propertyRequest.Latitude, propertyRequest.Longitude, propertyRequest.IsForSale, propertyRequest.SalePrice)
	if err != nil {
		return err
	}
	fmt.Println(property)

	return s.repo.RegisterNewProperty(property)
}

func (s *PropertyService) GetProperty(propertyName string) (*domain.Property, error) {
	return s.repo.GetProperty(propertyName)
}
