package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

// Initialise the database
func Init() error {
	var err error

	logger := log.Default()
	dbFile := "innbundet.sqlite"

	logger.Printf("Initialising database: %s", dbFile)
	db, err = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

// Get the database manager
func DbManager() *gorm.DB {
	return db
}
