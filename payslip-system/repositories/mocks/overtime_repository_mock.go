package mocks

import (
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockOvertimeRepo struct {
	mock.Mock
}

func (m *MockOvertimeRepo) OvertimeExists(userID uuid.UUID, date time.Time) bool {
	args := m.Called(userID, date)
	return args.Bool(0)
}

func (m *MockOvertimeRepo) CreateOvertime(userID uuid.UUID, hours int, date time.Time) error {
	args := m.Called(userID, hours, date)
	return args.Error(0)
}

func (m *MockOvertimeRepo) GetOvertimeHistory(userID uuid.UUID) ([]models.Overtime, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Overtime), args.Error(1)
}
