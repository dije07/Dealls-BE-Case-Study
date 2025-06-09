package interfaces

import (
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

type AttendanceRepository interface {
	AttendanceExists(userID uuid.UUID, date time.Time) bool
	CreateAttendance(userID uuid.UUID, date time.Time) error
	GetAttendanceHistory(userID uuid.UUID) ([]models.Attendance, error)
}
