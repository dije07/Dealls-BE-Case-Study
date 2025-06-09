package mocks

import (
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/stretchr/testify/mock"
)

type MockPayrollService struct {
	mock.Mock
}

func (m *MockPayrollService) CreatePayrollPeriod(start, end time.Time) error {
	args := m.Called(start, end)
	return args.Error(0)
}

func (m *MockPayrollService) RunPayroll(period models.PayrollPeriod) error {
	args := m.Called(period)
	return args.Error(0)
}
