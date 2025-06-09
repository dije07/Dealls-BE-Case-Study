package mocks

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPayslipRepo struct {
	mock.Mock
}

func (m *MockPayslipRepo) GetPayslip(userID, periodID uuid.UUID) (*models.Payslip, error) {
	args := m.Called(userID, periodID)
	p := args.Get(0)
	if p == nil {
		return nil, args.Error(1)
	}
	return p.(*models.Payslip), args.Error(1)
}

func (m *MockPayslipRepo) GetPayslipsByPeriod(periodID uuid.UUID) ([]models.Payslip, error) {
	args := m.Called(periodID)
	return args.Get(0).([]models.Payslip), args.Error(1)
}
