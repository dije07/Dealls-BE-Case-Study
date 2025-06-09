package mocks

import (
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPayrollRepo struct {
	mock.Mock
}

func (m *MockPayrollRepo) PayrollPeriodExists(start, end time.Time) bool {
	args := m.Called(start, end)
	return args.Bool(0)
}

func (m *MockPayrollRepo) CreatePayrollPeriod(start, end time.Time) error {
	args := m.Called(start, end)
	return args.Error(0)
}

func (m *MockPayrollRepo) GetAllEmployees() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockPayrollRepo) CountAttendances(userID uuid.UUID, start, end time.Time) (int, error) {
	args := m.Called(userID, start, end)
	return args.Int(0), args.Error(1)
}

func (m *MockPayrollRepo) SumOvertimeHours(userID uuid.UUID, start, end time.Time) (int, error) {
	args := m.Called(userID, start, end)
	return args.Int(0), args.Error(1)
}

func (m *MockPayrollRepo) SumReimbursements(userID uuid.UUID, start, end time.Time) (float64, error) {
	args := m.Called(userID, start, end)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockPayrollRepo) SavePayslip(p models.Payslip) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockPayrollRepo) ClosePayrollPeriod(periodID uuid.UUID) error {
	args := m.Called(periodID)
	return args.Error(0)
}
