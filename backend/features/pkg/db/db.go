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

	query := `
    CREATE TABLE IF NOT EXISTS properties (
    id uuid PRIMARY KEY NOT NULL, 
    name TEXT NOT NULL UNIQUE,                       
    neighborhood TEXT NOT NULL,              
    city TEXT NOT NULL,                     
    rent_amount INT NOT NULL CHECK (rent_amount >= 0),
    rent_currency CHAR(3) NOT NULL DEFAULT 'USD',    
    bedrooms INT NOT NULL CHECK (bedrooms >= 0),    
    bathrooms INT NOT NULL CHECK (bathrooms >= 0), 
    area NUMERIC(10, 2) NOT NULL CHECK (area >= 0),
    description TEXT,                         
    latitude NUMERIC(9, 6),                  
    longitude NUMERIC(9, 6),                
    created_at TIMESTAMP NOT NULL DEFAULT NOW(), 
    deleted_at TIMESTAMP,
    is_for_sale BOOLEAN NOT NULL DEFAULT FALSE,
    sale_price INT 
);

    `
	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	return db, nil
}
