package interfaces

import (
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

type OvertimeRepository interface {
	OvertimeExists(userID uuid.UUID, date time.Time) bool
	CreateOvertime(userID uuid.UUID, hours int, date time.Time) error
	GetOvertimeHistory(userID uuid.UUID) ([]models.Overtime, error)
}
