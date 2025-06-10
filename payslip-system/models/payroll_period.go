package models

import (
	"time"

	"github.com/google/uuid"
)

type PayrollPeriod struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	StartDate time.Time
	EndDate   time.Time
	IsClosed  bool
	IPAddress string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID `gorm:"type:uuid"`
	UpdatedBy uuid.UUID `gorm:"type:uuid"`
}
