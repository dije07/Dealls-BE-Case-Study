package services

import (
	"testing"
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/repositories/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSubmitOvertime_Success(t *testing.T) {
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	mockRepo := new(mocks.MockOvertimeRepo)
	mockRepo.On("OvertimeExists", userID, today).Return(false)
	mockRepo.On("CreateOvertime", userID, 2, today).Return(nil)

	service := &OvertimeServiceImpl{Repo: mockRepo}
	err := service.SubmitOvertime(userID, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubmitOvertime_TooManyHours(t *testing.T) {
	userID := uuid.New()
	service := &OvertimeServiceImpl{}

	err := service.SubmitOvertime(userID, 5) // > 3 hours
	assert.EqualError(t, err, "overtime must be between 1–3 hours")
}

func TestSubmitOvertime_AlreadySubmitted(t *testing.T) {
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	mockRepo := new(mocks.MockOvertimeRepo)
	mockRepo.On("OvertimeExists", userID, today).Return(true)

	service := &OvertimeServiceImpl{Repo: mockRepo}
	err := service.SubmitOvertime(userID, 2)

	assert.EqualError(t, err, "overtime already submitted for today")
	mockRepo.AssertExpectations(t)
}

func TestGetMyOvertime_Success(t *testing.T) {
	userID := uuid.New()

	mockRepo := new(mocks.MockOvertimeRepo)
	mockRepo.On("GetOvertimeHistory", userID).Return([]models.Overtime{
		{ID: uuid.New(), Hours: 3},
	}, nil)

	service := &OvertimeServiceImpl{Repo: mockRepo}
	result, err := service.GetMyOvertime(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

func TestNewOvertimeService_ReturnsImplementation(t *testing.T) {
	mockRepo := new(mocks.MockOvertimeRepo)
	service := NewOvertimeService(mockRepo)

	assert.NotNil(t, service)
	_, ok := service.(*OvertimeServiceImpl)
	assert.True(t, ok)
}
