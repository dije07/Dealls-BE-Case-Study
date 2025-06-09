package models

import (
	"time"

	"github.com/google/uuid"
)

type Overtime struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID
	User      User      `gorm:"foreignKey:UserID"`
	Date      time.Time `gorm:"index"` // Used for unique per-day constraint
	Hours     int
	CreatedAt time.Time
	UpdatedAt time.Time
}
