package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/FelipeMCassiano/buroka/features/property/domain"
)

var ErrNotKnownPropertyName = errors.New("property with this code does not matchs the name")

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

	query := `INSERT INTO properties (id,name, neighborhood, city, rent_amount, rent_currency, bedrooms, bathrooms, area, description, latitude, longitude, is_for_sale, sale_price, property_type, property_code, is_for_rent) 
                VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) 
    `

	_, err = tx.Exec(query, property.ID, property.Name, property.Neighborhood, property.City, property.Rent.Amount, property.Rent.Currency, property.Bedrooms, property.Bathrooms, property.Area, property.Description, property.Latitude, property.Longitude, property.IsForSale, property.SalePrice, property.PropertyType, property.Code, property.IsForRent)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return err
}

func (r *PropertyRepository) GetProperty(propertyName string, propertyCode string) (*domain.Property, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `SELECT id, name, neighborhood, city, rent_amount, rent_currency, bedrooms, bathrooms, area, description, latitude, longitude, is_for_sale, sale_price, created_at, deleted_at, property_type, property_code, is_for_rent FROM properties WHERE property_code=$1`

	property := new(domain.Property)
	if err := tx.QueryRow(query, propertyCode).Scan(&property.ID, &property.Name, &property.Neighborhood, &property.City, &property.Rent.Amount, &property.Rent.Currency, &property.Bedrooms, &property.Bathrooms, &property.Area, &property.Description, &property.Latitude, &property.Longitude, &property.IsForSale, &property.SalePrice, &property.CreatedAt, &property.DeletedAt, &property.PropertyType, &property.Code, &property.IsForRent); err != nil {
		return nil, err
	}

	if property.Name != propertyName {
		return nil, ErrNotKnownPropertyName
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return property, nil
}

func (r *PropertyRepository) SearchProperty(searchFilter domain.SearchFilter) ([]domain.Property, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	query := `

		SELECT id, name, neighborhood, city, rent_amount, rent_currency, bedrooms, bathrooms, area, description, 
		       latitude, longitude, is_for_sale, sale_price, created_at, deleted_at, property_type, 
		       property_code, is_for_rent
		FROM properties 
		WHERE 1=1
	`

	appendedQuery, args := appendConditions(query, searchFilter)
	fmt.Println(appendedQuery)

	rows, err := tx.Query(appendedQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	properties := []domain.Property{}
	for rows.Next() {
		var property domain.Property
		err := rows.Scan(
			&property.ID, &property.Name, &property.Neighborhood, &property.City,
			&property.Rent.Amount, &property.Rent.Currency, &property.Bedrooms, &property.Bathrooms,
			&property.Area, &property.Description, &property.Latitude, &property.Longitude,
			&property.IsForSale, &property.SalePrice, &property.CreatedAt, &property.DeletedAt,
			&property.PropertyType, &property.Code, &property.IsForRent,
		)
		if err != nil {
			return nil, err
		}
		properties = append(properties, property)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return properties, nil
}

func appendConditions(query string, searchFilter domain.SearchFilter) (string, []interface{}) {
	args := []interface{}{}
	argPosition := 1
	if searchFilter.PropertyType != "" {
		query += fmt.Sprintf(" AND property_type = $%d", argPosition)
		args = append(args, searchFilter.PropertyType)
		argPosition++
	}

	if searchFilter.ForRent {
		query += fmt.Sprintf(" AND is_for_rent = $%d", argPosition)
		args = append(args, true)
		argPosition++
	}
	if searchFilter.ForSale {
		query += fmt.Sprintf(" AND is_for_sale = $%d", argPosition)
		args = append(args, true)
		argPosition++
	}
	if searchFilter.Neighborhood != "" {
		query += fmt.Sprintf(" AND neighborhood ILIKE $%d", argPosition)
		args = append(args, "%"+searchFilter.Neighborhood+"%")
		argPosition++
	}

	if searchFilter.City != "" {
		query += fmt.Sprintf(" AND city ILIKE $%d", argPosition)
		args = append(args, "%"+searchFilter.City+"%")
		argPosition++
	}

	if searchFilter.RentAmount != 0 {
		query += fmt.Sprintf(" AND rent_amount = $%d", argPosition)
		args = append(args, searchFilter.RentAmount)
		argPosition++
	}

	if searchFilter.SalePrice != 0 {
		query += fmt.Sprintf(" AND sale_price = $%d", argPosition)
		args = append(args, searchFilter.SalePrice)
		argPosition++
	}

	if searchFilter.Size != 0 {
		query += fmt.Sprintf(" AND area = $%f", argPosition)
		args = append(args, searchFilter.Size)
		argPosition++
	}

	return query, args
}
