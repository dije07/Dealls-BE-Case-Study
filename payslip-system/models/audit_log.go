package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Action      string
	Entity      string
	EntityID    uuid.UUID
	PerformedBy uuid.UUID
	IP          string
	RequestID   string
	CreatedAt   time.Time
}
