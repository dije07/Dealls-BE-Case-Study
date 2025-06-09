package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/services/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSubmitReimbursement_Success(t *testing.T) {
	e := echo.New()
	body := map[string]interface{}{
		"amount":      150000,
		"description": "Transport reimbursement",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/reimbursement", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userID := uuid.New()
	c.Set("user_id", userID.String())

	mockService := new(mocks.MockReimbursementService)
	mockService.On("SubmitReimbursement", userID, 150000.0, "Transport reimbursement").Return(nil)

	handler := &ReimbursementHandler{Service: mockService}
	err := handler.SubmitReimbursement(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "reimbursement submitted")
	mockService.AssertExpectations(t)
}

func TestSubmitReimbursement_InvalidInput(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/reimbursement", bytes.NewReader([]byte("invalid")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user_id", uuid.New().String())

	handler := &ReimbursementHandler{Service: new(mocks.MockReimbursementService)}
	err := handler.SubmitReimbursement(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid input")
}

func TestGetMyReimbursements_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/reimbursement", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	userID := uuid.New()
	c.Set("user_id", userID.String())

	mockService := new(mocks.MockReimbursementService)
	mockService.On("GetMyReimbursements", userID).Return([]models.Reimbursement{
		{
			ID:          uuid.New(),
			Amount:      150000,
			Description: "Lunch",
			CreatedAt:   time.Now(),
		},
	}, nil)

	handler := &ReimbursementHandler{Service: mockService}
	err := handler.GetMyReimbursements(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "reimbursements")
	mockService.AssertExpectations(t)
}

func TestGetMyReimbursements_Unauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/reimbursement", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &ReimbursementHandler{Service: new(mocks.MockReimbursementService)}
	err := handler.GetMyReimbursements(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestSubmitReimbursement_Unauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/reimbursement", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &ReimbursementHandler{Service: new(mocks.MockReimbursementService)}
	err := handler.SubmitReimbursement(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestSubmitReimbursement_ServiceError(t *testing.T) {
	e := echo.New()
	body := map[string]interface{}{
		"amount":      0, // triggers error from service
		"description": "testing",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/reimbursement", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	c.Set("user_id", testUserID.String())

	mockService := new(mocks.MockReimbursementService)
	mockService.On("SubmitReimbursement", testUserID, 0.0, "testing").
		Return(errors.New("amount must be greater than zero"))

	handler := &ReimbursementHandler{Service: mockService}
	err := handler.SubmitReimbursement(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "amount must be greater than zero")
	mockService.AssertExpectations(t)
}

func TestGetMyReimbursements_FetchError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/reimbursement", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	c.Set("user_id", testUserID.String())

	mockService := new(mocks.MockReimbursementService)
	mockService.On("GetMyReimbursements", testUserID).
		Return(nil, errors.New("db error"))

	handler := &ReimbursementHandler{Service: mockService}
	err := handler.GetMyReimbursements(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to fetch reimbursements")
	mockService.AssertExpectations(t)
}
