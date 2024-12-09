package application

import (
	"github.com/FelipeMCassiano/buroka/features/property/domain"
	"github.com/FelipeMCassiano/buroka/features/property/infrastructure"
)

type PropertyService struct {
	repo *infrastructure.PropertyRepository
}

type RegisterPropertyRequest struct {
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
	PropertyType string       `json:"property_type"`
	IsForRent    bool         `json:"is_for_rent"`
}

func NewPropertyService(repo *infrastructure.PropertyRepository) *PropertyService {
	return &PropertyService{repo: repo}
}

func (s *PropertyService) RegisterProperty(propertyRequest *RegisterPropertyRequest) (*domain.Property, error) {
	property, err := domain.NewProperty(propertyRequest.Neighborhood, propertyRequest.City, propertyRequest.Description, propertyRequest.Rent, propertyRequest.Bedrooms, propertyRequest.Bathrooms, propertyRequest.Area, propertyRequest.Latitude, propertyRequest.Longitude, propertyRequest.IsForSale, propertyRequest.SalePrice, propertyRequest.PropertyType, propertyRequest.IsForRent)
	if err != nil {
		return nil, err
	}
	if err := s.repo.RegisterNewProperty(property); err != nil {
		return nil, err
	}

	return property, nil
}

func (s *PropertyService) GetProperty(propertyName string, propertyCode string) (*domain.Property, error) {
	return s.repo.GetProperty(propertyName, propertyCode)
}

func (s *PropertyService) SearchProperty(searchFilter domain.SearchFilter) ([]domain.Property, error) {
	return s.repo.SearchProperty(searchFilter)
}
