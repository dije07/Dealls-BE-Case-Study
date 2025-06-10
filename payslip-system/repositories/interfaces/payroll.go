package interfaces

import (
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PayrollRepository interface {
	PayrollPeriodExists(start, end time.Time) bool
	CreatePayrollPeriod(c echo.Context, userID uuid.UUID, start, end time.Time) error
	GetAllEmployees() ([]models.User, error)
	CountAttendances(userID uuid.UUID, start, end time.Time) (int, error)
	SumOvertimeHours(userID uuid.UUID, start, end time.Time) (int, error)
	SumReimbursements(userID uuid.UUID, start, end time.Time) (float64, error)
	SavePayslip(p models.Payslip) error
	ClosePayrollPeriod(periodID uuid.UUID) error
}
