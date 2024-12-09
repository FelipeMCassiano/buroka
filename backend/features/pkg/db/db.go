package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return nil, fmt.Errorf("DATABASE URL enviroment variable empty")
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		return nil, err
	}

	query1 := `
    CREATE EXTENSION IF NOT EXISTS pg_trgm;
    `
	if _, err := db.Exec(query1); err != nil {
		return nil, err
	}

	query2 := `
    CREATE TABLE IF NOT EXISTS properties (
    id uuid PRIMARY KEY NOT NULL, 
    name TEXT NOT NULL,                       
    property_type TEXT NOT NULL,
    neighborhood TEXT NOT NULL,              
    city TEXT NOT NULL,                     
    rent_amount INT CHECK (rent_amount >= 0),
    rent_currency CHAR(3) DEFAULT 'USD',    
    bedrooms INT NOT NULL CHECK (bedrooms >= 0),    
    bathrooms INT NOT NULL CHECK (bathrooms >= 0), 
    area NUMERIC(10, 2) NOT NULL CHECK (area >= 0),
    description TEXT,                         
    latitude NUMERIC(9, 6),                  
    longitude NUMERIC(9, 6),                
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), 
    deleted_at TIMESTAMP,
    is_for_sale BOOLEAN NOT NULL DEFAULT FALSE,
    is_for_rent BOOLEAN NOT NULL DEFAULT FALSE,
    sale_price INT,
    property_code VARCHAR(10) NOT NULL UNIQUE
    );
    CREATE INDEX IF NOT EXISTS idx_neighborhood_trgm ON properties USING gin (neighborhood gin_trgm_ops);
    `

	if _, err := db.Exec(query2); err != nil {
		return nil, err
	}

	return db, nil
}
