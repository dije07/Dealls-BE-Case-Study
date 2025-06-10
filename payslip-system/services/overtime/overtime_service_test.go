package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/repositories/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSubmitOvertime_Success(t *testing.T) {
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	mockRepo := new(mocks.MockOvertimeRepo)
	mockRepo.On("OvertimeExists", userID, today).Return(false)
	mockRepo.On("CreateOvertime", mock.AnythingOfType("*echo.context"), userID, 2, mock.AnythingOfType("time.Time")).Return(nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	service := &OvertimeServiceImpl{Repo: mockRepo}
	err := service.SubmitOvertime(c, userID, 2)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubmitOvertime_TooManyHours(t *testing.T) {
	userID := uuid.New()
	service := &OvertimeServiceImpl{}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := service.SubmitOvertime(c, userID, 5) // > 3 hours
	assert.EqualError(t, err, "overtime must be between 1–3 hours")
}

func TestSubmitOvertime_AlreadySubmitted(t *testing.T) {
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	mockRepo := new(mocks.MockOvertimeRepo)
	mockRepo.On("OvertimeExists", userID, today).Return(true)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	service := &OvertimeServiceImpl{Repo: mockRepo}
	err := service.SubmitOvertime(c, userID, 2)

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
