package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	Username  string    `gorm:"uniqueIndex"`
	Password  string
	RoleID    uint
	Role      Role    `gorm:"foreignKey:RoleID"`
	Salary    float64 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
