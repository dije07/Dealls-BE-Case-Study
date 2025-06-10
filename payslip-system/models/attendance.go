package models

import (
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID
	User      User      `gorm:"foreignKey:UserID"`
	Date      time.Time `gorm:"index"`
	IPAddress string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID `gorm:"type:uuid"`
	UpdatedBy uuid.UUID `gorm:"type:uuid"`
}
