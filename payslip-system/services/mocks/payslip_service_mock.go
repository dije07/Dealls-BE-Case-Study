package mocks

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPayslipService struct {
	mock.Mock
}

func (m *MockPayslipService) GetEmployeePayslip(userID, periodID uuid.UUID) (*models.Payslip, error) {
	args := m.Called(userID, periodID)
	p := args.Get(0)
	if p == nil {
		return nil, args.Error(1)
	}
	return p.(*models.Payslip), args.Error(1)
}

func (m *MockPayslipService) GetPayslipSummary(periodID uuid.UUID) ([]models.Payslip, float64, error) {
	args := m.Called(periodID)
	payslips := args.Get(0).([]models.Payslip)
	total := args.Get(1).(float64)
	return payslips, total, args.Error(2)
}
