package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/services/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetMyPayslip_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslip/:period_id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	testPeriodID := uuid.New()

	c.Set("user_id", testUserID.String())
	c.SetParamNames("period_id")
	c.SetParamValues(testPeriodID.String())

	mockService := new(mocks.MockPayslipService)
	mockService.On("GetEmployeePayslip", testUserID, testPeriodID).Return(&models.Payslip{
		UserID:         testUserID,
		PeriodID:       testPeriodID,
		BaseSalary:     8000000,
		AttendanceDays: 20,
		OvertimeHours:  3,
		OvertimePay:    400000,
		Reimbursement:  150000,
		TakeHomePay:    8550000,
	}, nil)

	handler := &PayslipHandler{Service: mockService}
	err := handler.GetMyPayslip(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "take_home_pay")
	mockService.AssertExpectations(t)
}

func TestGetMyPayslip_InvalidUUID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslip/:period_id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Set("user_id", uuid.New().String())
	c.SetParamNames("period_id")
	c.SetParamValues("not-a-uuid")

	handler := &PayslipHandler{}
	err := handler.GetMyPayslip(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid payroll period ID")
}

func TestGetMyPayslip_NotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslip/:period_id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testUserID := uuid.New()
	testPeriodID := uuid.New()

	c.Set("user_id", testUserID.String())
	c.SetParamNames("period_id")
	c.SetParamValues(testPeriodID.String())

	mockService := new(mocks.MockPayslipService)
	mockService.On("GetEmployeePayslip", testUserID, testPeriodID).Return(nil, errors.New("not found"))

	handler := &PayslipHandler{Service: mockService}
	err := handler.GetMyPayslip(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "payslip not found")
	mockService.AssertExpectations(t)
}

func TestGetPayslipSummary_Success(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslip-summary/:period_id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testPeriodID := uuid.New()
	c.SetParamNames("period_id")
	c.SetParamValues(testPeriodID.String())

	mockService := new(mocks.MockPayslipService)
	mockService.On("GetPayslipSummary", testPeriodID).Return([]models.Payslip{
		{
			UserID:         uuid.New(),
			BaseSalary:     10000000,
			AttendanceDays: 20,
			OvertimeHours:  2,
			OvertimePay:    500000,
			Reimbursement:  200000,
			TakeHomePay:    10700000,
		},
	}, 10700000.0, nil)

	handler := &PayslipHandler{Service: mockService}
	err := handler.GetPayslipSummary(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "summary")
	assert.Contains(t, rec.Body.String(), "total_take_home")
	mockService.AssertExpectations(t)
}

func TestGetPayslipSummary_InvalidUUID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslip-summary/:period_id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("period_id")
	c.SetParamValues("not-a-uuid")

	handler := &PayslipHandler{}
	err := handler.GetPayslipSummary(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid period ID")
}

func TestGetPayslipSummary_ServiceError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslip-summary/:period_id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	testPeriodID := uuid.New()
	c.SetParamNames("period_id")
	c.SetParamValues(testPeriodID.String())

	mockService := new(mocks.MockPayslipService)
	mockService.On("GetPayslipSummary", testPeriodID).Return([]models.Payslip{}, 0.0, errors.New("fetch error"))

	handler := &PayslipHandler{Service: mockService}
	err := handler.GetPayslipSummary(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to fetch summary")
	mockService.AssertExpectations(t)
}

func TestGetMyPayslip_Unauthorized(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/payslip/:period_id", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// no user_id set
	c.SetParamNames("period_id")
	c.SetParamValues(uuid.New().String())

	handler := &PayslipHandler{}
	err := handler.GetMyPayslip(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}
