package db

import (
	"github.com/michaelenger/innbundet/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initialise the database
func Init(filePath string) (*gorm.DB, error) {
	log.Debug("Initialising database: %s", filePath)
	db, err := gorm.Open(sqlite.Open(filePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
