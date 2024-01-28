package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Initialise the database
func Init(filePath string) (*gorm.DB, error) {
	logger := log.Default()

	logger.Printf("Initialising database: %s", filePath)
	db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
