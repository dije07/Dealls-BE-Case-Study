package interfaces

import (
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AttendanceRepository interface {
	AttendanceExists(userID uuid.UUID, date time.Time) bool
	CreateAttendance(c echo.Context, userID uuid.UUID, date time.Time) error
	GetAttendanceHistory(userID uuid.UUID) ([]models.Attendance, error)
}
