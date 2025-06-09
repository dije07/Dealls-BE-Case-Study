package models

import "time"

type Role struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
