package interfaces

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

type ReimbursementRepository interface {
	CreateReimbursement(userID uuid.UUID, amount float64, description string) error
	GetReimbursementsByUser(userID uuid.UUID) ([]models.Reimbursement, error)
}
