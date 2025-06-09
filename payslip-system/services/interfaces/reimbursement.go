package interfaces

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

type ReimbursementService interface {
	SubmitReimbursement(userID uuid.UUID, amount float64, description string) error
	GetMyReimbursements(userID uuid.UUID) ([]models.Reimbursement, error)
}
