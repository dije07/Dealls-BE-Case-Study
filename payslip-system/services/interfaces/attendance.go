package interfaces

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AttendanceService interface {
	SubmitAttendance(c echo.Context, userID uuid.UUID) error
	GetMyAttendance(userID uuid.UUID) ([]models.Attendance, error)
}
