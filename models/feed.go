package models

import (
	"time"
)

// A feed that is being subscribed to
type Feed struct {
	ID          uint `gorm:"primaryKey"`
	Url         string
	Title       string
	Link        string
	Description string
	Image       *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
