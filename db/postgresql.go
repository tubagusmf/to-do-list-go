package db

import (
	"log"
	"to-do-list/internal/helper"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres initializes a PostgreSQL database connection using Gorm.
func NewPostgres() *gorm.DB {
	// Get connection string
	dsn := helper.GetConnectionString()

	// Connect to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}
