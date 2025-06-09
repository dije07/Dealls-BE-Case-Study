package services

import (
	"github.com/dije07/payslip-system/models"
	repositoryInterfaces "github.com/dije07/payslip-system/repositories/interfaces"
	"github.com/dije07/payslip-system/services/interfaces"
	"github.com/google/uuid"
)

type PayslipServiceImpl struct {
	Repo repositoryInterfaces.PayslipRepository
}

func NewPayslipService(repo repositoryInterfaces.PayslipRepository) interfaces.PayslipService {
	return &PayslipServiceImpl{Repo: repo}
}

func (s *PayslipServiceImpl) GetEmployeePayslip(userID, periodID uuid.UUID) (*models.Payslip, error) {
	return s.Repo.GetPayslip(userID, periodID)
}

func (s *PayslipServiceImpl) GetPayslipSummary(periodID uuid.UUID) ([]models.Payslip, float64, error) {
	payslips, err := s.Repo.GetPayslipsByPeriod(periodID)
	if err != nil {
		return nil, 0, err
	}

	var total float64
	for _, p := range payslips {
		total += p.TakeHomePay
	}

	return payslips, total, nil
}
