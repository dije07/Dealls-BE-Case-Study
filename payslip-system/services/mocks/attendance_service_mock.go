package mocks

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

type MockAttendanceService struct {
	mock.Mock
}

func (m *MockAttendanceService) SubmitAttendance(c echo.Context, userID uuid.UUID) error {
	args := m.Called(c, userID)
	return args.Error(0)
}

func (m *MockAttendanceService) GetMyAttendance(userID uuid.UUID) ([]models.Attendance, error) {
	args := m.Called(userID)
	// Safely cast to the expected return types
	result, _ := args.Get(0).([]models.Attendance)
	return result, args.Error(1)
}
