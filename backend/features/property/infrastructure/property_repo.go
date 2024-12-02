package infrastructure

import (
	"database/sql"

	"github.com/FelipeMCassiano/buroka/features/property/domain"
)

type PropertyRepository struct {
	db *sql.DB
}

func NewPropertyRepository(db *sql.DB) *PropertyRepository {
	return &PropertyRepository{db: db}
}

func (r *PropertyRepository) RegisterNewProperty(property *domain.Property) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `INSERT INTO properties (id,name, neighborhood, city, rent_amount, rent_currency, bedrooms, bathrooms, area, description, latitude, longitude, is_for_sale, sale_price) 
                VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) 
    `

	_, err = tx.Exec(query, property.ID, property.Name, property.Neighborhood, property.City, property.Rent.Amount, property.Rent.Currency, property.Bedrooms, property.Bathrooms, property.Area, property.Description, property.Latitude, property.Longitude, property.IsForSale, property.SalePrice)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return err
}

func (r *PropertyRepository) GetProperty(propertyName string) (*domain.Property, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `SELECT id, name, neighborhood, city, rent_amount, rent_currency, bedrooms, bathrooms, area, description, latitude, longitude, is_for_sale, sale_price, created_at, deleted_at FROM properties WHERE name=$1`

	property := new(domain.Property)
	if err := tx.QueryRow(query, propertyName).Scan(&property.ID, &property.Name, &property.Neighborhood, &property.City, &property.Rent.Amount, &property.Rent.Currency, &property.Bedrooms, &property.Bathrooms, &property.Area, &property.Description, &property.Latitude, &property.Longitude, &property.IsForSale, &property.SalePrice, &property.CreatedAt, &property.DeletedAt); err != nil {
		return nil, err
	}

	return property, nil
}
