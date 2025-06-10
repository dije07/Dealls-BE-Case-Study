package services

import (
	"errors"

	"github.com/dije07/payslip-system/models"
	repositoryInterfaces "github.com/dije07/payslip-system/repositories/interfaces"
	"github.com/dije07/payslip-system/services/interfaces"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReimbursementServiceImpl struct {
	Repo repositoryInterfaces.ReimbursementRepository
}

func NewReimbursementService(repo repositoryInterfaces.ReimbursementRepository) interfaces.ReimbursementService {
	return &ReimbursementServiceImpl{Repo: repo}
}

func (s *ReimbursementServiceImpl) SubmitReimbursement(c echo.Context, userID uuid.UUID, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if len(description) == 0 {
		return errors.New("description is required")
	}
	return s.Repo.CreateReimbursement(c, userID, amount, description)
}

func (s *ReimbursementServiceImpl) GetMyReimbursements(userID uuid.UUID) ([]models.Reimbursement, error) {
	return s.Repo.GetReimbursementsByUser(userID)
}
