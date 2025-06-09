package interfaces

import (
	"time"

	"github.com/dije07/payslip-system/models"
)

type PayrollService interface {
	CreatePayrollPeriod(start, end time.Time) error
	RunPayroll(period models.PayrollPeriod) error
}
