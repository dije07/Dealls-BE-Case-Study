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
	CreatedAt time.Time
	UpdatedAt time.Time
}
