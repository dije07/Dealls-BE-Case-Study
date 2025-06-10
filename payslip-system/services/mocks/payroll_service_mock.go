package mocks

import (
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

type MockPayrollService struct {
	mock.Mock
}

func (m *MockPayrollService) CreatePayrollPeriod(c echo.Context, userID uuid.UUID, start, end time.Time) error {
	args := m.Called(c, userID, start, end)
	return args.Error(0)
}

func (m *MockPayrollService) RunPayroll(period models.PayrollPeriod) error {
	args := m.Called(period)
	return args.Error(0)
}
