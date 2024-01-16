package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// Initialise the database
func Init() (*gorm.DB, error) {
	var err error

	logger := log.Default()
	dbFile := "innbundet.sqlite"

	logger.Printf("Initialising database: %s", dbFile)
	db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Get the database manager
func DbManager() *gorm.DB {
	return db
}
