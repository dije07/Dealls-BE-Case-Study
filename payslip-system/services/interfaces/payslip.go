package interfaces

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

type PayslipService interface {
	GetEmployeePayslip(userID, periodID uuid.UUID) (*models.Payslip, error)
	GetPayslipSummary(periodID uuid.UUID) ([]models.Payslip, float64, error)
}
