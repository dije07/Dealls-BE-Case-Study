package handlers

import (
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

func TestSubmitAttendance_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/attendance", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	c.Set("user_id", testUserID.String())

	mockService := new(mocks.MockAttendanceService)
	mockService.On("SubmitAttendance", testUserID).Return(nil)

	handler := &AttendanceHandler{Service: mockService}
	err := handler.SubmitAttendance(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "attendance submitted")
	mockService.AssertExpectations(t)
}

func TestSubmitAttendance_Unauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/attendance", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// No user_id set in context
	mockService := new(mocks.MockAttendanceService)
	handler := &AttendanceHandler{Service: mockService}

	err := handler.SubmitAttendance(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestSubmitAttendance_AlreadySubmitted(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/attendance", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	c.Set("user_id", testUserID.String())

	mockService := new(mocks.MockAttendanceService)
	mockService.On("SubmitAttendance", testUserID).Return(errors.New("attendance already submitted for today"))

	handler := &AttendanceHandler{Service: mockService}
	err := handler.SubmitAttendance(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "attendance already submitted")
	mockService.AssertExpectations(t)
}

func TestGetMyAttendance_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/attendance", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	c.Set("user_id", testUserID.String())

	mockService := new(mocks.MockAttendanceService)
	mockService.On("GetMyAttendance", testUserID).Return([]models.Attendance{
		{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
		},
	}, nil)

	handler := &AttendanceHandler{Service: mockService}
	err := handler.GetMyAttendance(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "attendance")
	mockService.AssertExpectations(t)
}

func TestGetMyAttendance_Unauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/attendance", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// no user_id in context
	mockService := new(mocks.MockAttendanceService)
	handler := &AttendanceHandler{Service: mockService}

	err := handler.GetMyAttendance(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}

func TestGetMyAttendance_ErrorFromService(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/attendance", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	c.Set("user_id", testUserID.String())

	mockService := new(mocks.MockAttendanceService)
	mockService.On("GetMyAttendance", testUserID).Return(nil, errors.New("database failure"))

	handler := &AttendanceHandler{Service: mockService}
	err := handler.GetMyAttendance(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to fetch attendance")
	mockService.AssertExpectations(t)
}
