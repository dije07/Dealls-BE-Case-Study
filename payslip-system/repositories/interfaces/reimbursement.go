package interfaces

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReimbursementRepository interface {
	CreateReimbursement(c echo.Context, userID uuid.UUID, amount float64, description string) error
	GetReimbursementsByUser(userID uuid.UUID) ([]models.Reimbursement, error)
}
