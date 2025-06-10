package interfaces

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReimbursementService interface {
	SubmitReimbursement(c echo.Context, userID uuid.UUID, amount float64, description string) error
	GetMyReimbursements(userID uuid.UUID) ([]models.Reimbursement, error)
}
