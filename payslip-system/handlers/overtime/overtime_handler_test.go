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

func TestSubmitOvertime_Success(t *testing.T) {
	e := echo.New()
	body := map[string]interface{}{"hours": 2}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/overtime", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	c.Set("user_id", testUserID.String())
	c.Set("role", "employee")

	mockService := new(mocks.MockOvertimeService)
	mockService.On("SubmitOvertime", c, testUserID, 2).Return(nil)

	handler := &OvertimeHandler{Service: mockService}
	err := handler.SubmitOvertime(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "overtime submitted")

	mockService.AssertExpectations(t)
}

func TestSubmitOvertime_Unauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/overtime", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(mocks.MockOvertimeService)
	handler := &OvertimeHandler{Service: mockService}

	err := handler.SubmitOvertime(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestSubmitOvertime_InvalidInput(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/overtime", bytes.NewReader([]byte("invalid")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	c.Set("user_id", testUserID.String())

	mockService := new(mocks.MockOvertimeService)
	handler := &OvertimeHandler{Service: mockService}

	err := handler.SubmitOvertime(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid input")
}

func TestSubmitOvertime_AlreadySubmitted(t *testing.T) {
	e := echo.New()
	body := map[string]interface{}{"hours": 3}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/overtime", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	c.Set("user_id", testUserID.String())

	mockService := new(mocks.MockOvertimeService)
	mockService.On("SubmitOvertime", c, testUserID, 3).Return(errors.New("overtime already submitted"))

	handler := &OvertimeHandler{Service: mockService}
	err := handler.SubmitOvertime(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "overtime already submitted")
	mockService.AssertExpectations(t)
}

func TestGetMyOvertime_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/overtime", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	c.Set("user_id", testUserID.String())

	mockService := new(mocks.MockOvertimeService)
	mockService.On("GetMyOvertime", testUserID).Return([]models.Overtime{
		{
			ID:        uuid.New(),
			Date:      time.Now(),
			Hours:     2,
			CreatedAt: time.Now(),
		},
	}, nil)

	handler := &OvertimeHandler{Service: mockService}
	err := handler.GetMyOvertime(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "overtime")
	mockService.AssertExpectations(t)
}

func TestGetMyOvertime_Unauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/overtime", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(mocks.MockOvertimeService)
	handler := &OvertimeHandler{Service: mockService}

	err := handler.GetMyOvertime(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestGetMyOvertime_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/overtime", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	c.Set("user_id", testUserID.String())

	mockService := new(mocks.MockOvertimeService)
	mockService.On("GetMyOvertime", testUserID).Return(nil, errors.New("DB error"))

	handler := &OvertimeHandler{Service: mockService}
	err := handler.GetMyOvertime(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to fetch overtime history")
	mockService.AssertExpectations(t)
}
