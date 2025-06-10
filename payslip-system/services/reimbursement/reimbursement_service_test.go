package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/repositories/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSubmitReimbursement_Success(t *testing.T) {
	userID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	mockRepo := new(mocks.MockReimbursementRepo)
	mockRepo.On("CreateReimbursement", c, userID, 100000.0, "Lunch").Return(nil)

	service := NewReimbursementService(mockRepo)
	err := service.SubmitReimbursement(c, userID, 100000, "Lunch")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubmitReimbursement_ZeroAmount(t *testing.T) {
	userID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	service := NewReimbursementService(new(mocks.MockReimbursementRepo))
	err := service.SubmitReimbursement(c, userID, 0, "Food")

	assert.EqualError(t, err, "amount must be greater than zero")
}

func TestSubmitReimbursement_EmptyDescription(t *testing.T) {
	userID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	service := NewReimbursementService(new(mocks.MockReimbursementRepo))
	err := service.SubmitReimbursement(c, userID, 150000, "")

	assert.EqualError(t, err, "description is required")
}

func TestGetMyReimbursements_Success(t *testing.T) {
	userID := uuid.New()
	mockData := []models.Reimbursement{
		{ID: uuid.New(), Amount: 120000, Description: "Taxi"},
	}

	mockRepo := new(mocks.MockReimbursementRepo)
	mockRepo.On("GetReimbursementsByUser", userID).Return(mockData, nil)

	service := NewReimbursementService(mockRepo)
	result, err := service.GetMyReimbursements(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Taxi", result[0].Description)
	mockRepo.AssertExpectations(t)
}

func TestNewReimbursementService_ReturnsImpl(t *testing.T) {
	mockRepo := new(mocks.MockReimbursementRepo)
	service := NewReimbursementService(mockRepo)

	assert.NotNil(t, service)
	_, ok := service.(*ReimbursementServiceImpl)
	assert.True(t, ok)
}
