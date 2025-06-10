package mocks

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
)

type MockReimbursementRepo struct {
	mock.Mock
}

func (m *MockReimbursementRepo) CreateReimbursement(c echo.Context, userID uuid.UUID, amount float64, description string) error {
	args := m.Called(c, userID, amount, description)
	return args.Error(0)
}

func (m *MockReimbursementRepo) GetReimbursementsByUser(userID uuid.UUID) ([]models.Reimbursement, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Reimbursement), args.Error(1)
}
