package mocks

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

type MockOvertimeService struct {
	mock.Mock
}

func (m *MockOvertimeService) SubmitOvertime(c echo.Context, userID uuid.UUID, hours int) error {
	args := m.Called(c, userID, hours)
	return args.Error(0)
}

func (m *MockOvertimeService) GetMyOvertime(userID uuid.UUID) ([]models.Overtime, error) {
	args := m.Called(userID)
	// Safely cast to the expected return types
	result, _ := args.Get(0).([]models.Overtime)
	return result, args.Error(1)
}
