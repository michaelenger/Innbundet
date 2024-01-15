package models

import (
	"database/sql"
	"time"
)

// A feed that is being subscribed to
type Feed struct {
	ID          uint `gorm:"primaryKey"`
	Url         string
	Title       string
	Link        sql.NullString
	Description sql.NullString
	Image       sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
