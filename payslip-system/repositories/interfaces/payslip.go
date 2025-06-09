package interfaces

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

type PayslipRepository interface {
	GetPayslip(userID, periodID uuid.UUID) (*models.Payslip, error)
	GetPayslipsByPeriod(periodID uuid.UUID) ([]models.Payslip, error)
}
