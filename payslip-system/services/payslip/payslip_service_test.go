package services

import (
	"errors"
	"testing"

	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/repositories/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetEmployeePayslip_Success(t *testing.T) {
	userID := uuid.New()
	periodID := uuid.New()
	expected := &models.Payslip{
		UserID:        userID,
		PeriodID:      periodID,
		TakeHomePay:   10000000,
		OvertimeHours: 2,
	}

	mockRepo := new(mocks.MockPayslipRepo)
	mockRepo.On("GetPayslip", userID, periodID).Return(expected, nil)

	service := NewPayslipService(mockRepo)
	result, err := service.GetEmployeePayslip(userID, periodID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetEmployeePayslip_Error(t *testing.T) {
	userID := uuid.New()
	periodID := uuid.New()

	mockRepo := new(mocks.MockPayslipRepo)
	mockRepo.On("GetPayslip", userID, periodID).Return(nil, errors.New("not found"))

	service := NewPayslipService(mockRepo)
	result, err := service.GetEmployeePayslip(userID, periodID)

	assert.Nil(t, result)
	assert.EqualError(t, err, "not found")
	mockRepo.AssertExpectations(t)
}

func TestGetPayslipSummary_Success(t *testing.T) {
	periodID := uuid.New()
	payslips := []models.Payslip{
		{UserID: uuid.New(), TakeHomePay: 5000000},
		{UserID: uuid.New(), TakeHomePay: 4000000},
	}

	mockRepo := new(mocks.MockPayslipRepo)
	mockRepo.On("GetPayslipsByPeriod", periodID).Return(payslips, nil)

	service := NewPayslipService(mockRepo)
	result, total, err := service.GetPayslipSummary(periodID)

	assert.NoError(t, err)
	assert.Equal(t, payslips, result)
	assert.Equal(t, 9000000.0, total)
	mockRepo.AssertExpectations(t)
}

func TestGetPayslipSummary_Error(t *testing.T) {
	periodID := uuid.New()

	mockRepo := new(mocks.MockPayslipRepo)
	mockRepo.On("GetPayslipsByPeriod", periodID).Return([]models.Payslip{}, errors.New("db error"))

	service := NewPayslipService(mockRepo)
	result, total, err := service.GetPayslipSummary(periodID)

	assert.Nil(t, result)
	assert.Equal(t, 0.0, total)
	assert.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

func TestNewPayslipService_ReturnsImpl(t *testing.T) {
	mockRepo := new(mocks.MockPayslipRepo)
	service := NewPayslipService(mockRepo)

	assert.NotNil(t, service)
	_, ok := service.(*PayslipServiceImpl)
	assert.True(t, ok)
}
