package models

import (
	"time"

	"github.com/google/uuid"
)

type Reimbursement struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID
	User        User `gorm:"foreignKey:UserID"`
	Amount      float64
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
