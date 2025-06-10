package interfaces

import (
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PayrollService interface {
	CreatePayrollPeriod(c echo.Context, userID uuid.UUID, start, end time.Time) error
	RunPayroll(period models.PayrollPeriod) error
}
