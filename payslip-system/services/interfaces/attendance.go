package interfaces

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

type AttendanceService interface {
	SubmitAttendance(userID uuid.UUID) error
	GetMyAttendance(userID uuid.UUID) ([]models.Attendance, error)
}
