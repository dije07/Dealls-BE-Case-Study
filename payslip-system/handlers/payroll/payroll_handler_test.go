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

func TestCreatePayrollPeriod_Success(t *testing.T) {
	e := echo.New()

	body := map[string]string{
		"start_date": "2025-06-01",
		"end_date":   "2025-06-30",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/payroll-period", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	id := uuid.New()
	c.Set("user_id", id.String())

	// Mock service and handler
	mockService := new(mocks.MockPayrollService)
	start, _ := time.Parse("2006-01-02", "2025-06-01")
	end, _ := time.Parse("2006-01-02", "2025-06-30")
	mockService.On("CreatePayrollPeriod", c, id, start, end).Return(nil)

	handler := &PayrollHandler{Service: mockService}

	err := handler.CreatePayrollPeriod(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "payroll period created")
}

func TestCreatePayrollPeriod_InvalidInput(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/payroll-period", bytes.NewReader([]byte("invalid")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &PayrollHandler{}
	err := handler.CreatePayrollPeriod(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid input")
}

func TestCreatePayrollPeriod_InvalidDateFormat(t *testing.T) {
	e := echo.New()
	body := map[string]string{
		"start_date": "June 1, 2025",
		"end_date":   "June 30, 2025",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/payroll-period", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &PayrollHandler{}
	err := handler.CreatePayrollPeriod(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid date format")
}

func TestCreatePayrollPeriod_ServiceError(t *testing.T) {
	e := echo.New()
	body := map[string]string{
		"start_date": "2025-06-01",
		"end_date":   "2025-06-30",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/payroll-period", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	id := uuid.New()
	c.Set("user_id", id.String())

	start, _ := time.Parse("2006-01-02", body["start_date"])
	end, _ := time.Parse("2006-01-02", body["end_date"])

	mockService := new(mocks.MockPayrollService)
	mockService.On("CreatePayrollPeriod", c, id, start, end).Return(errors.New("period exists"))

	handler := &PayrollHandler{
		Service: mockService,
		GetPeriodByID: func(id uuid.UUID) (models.PayrollPeriod, error) {
			return models.PayrollPeriod{ID: id}, nil
		},
	}

	err := handler.CreatePayrollPeriod(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "period exists")
	mockService.AssertExpectations(t)
}

func TestRunPayroll_Success(t *testing.T) {
	e := echo.New()
	periodID := uuid.New()
	body := map[string]string{"period_id": periodID.String()}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/run-payroll", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(mocks.MockPayrollService)
	mockService.On("RunPayroll", models.PayrollPeriod{ID: periodID}).Return(nil)

	handler := &PayrollHandler{
		Service: mockService,
		GetPeriodByID: func(id uuid.UUID) (models.PayrollPeriod, error) {
			return models.PayrollPeriod{ID: id}, nil
		},
	}

	err := handler.RunPayroll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "payroll processed")
	mockService.AssertExpectations(t)
}

func TestRunPayroll_InvalidInput(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/run-payroll", bytes.NewReader([]byte("bad json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &PayrollHandler{}
	err := handler.RunPayroll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid input")
}

func TestRunPayroll_PeriodNotFound(t *testing.T) {
	e := echo.New()
	body := map[string]string{"period_id": uuid.New().String()}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/run-payroll", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(mocks.MockPayrollService)

	handler := &PayrollHandler{
		Service: mockService,
		GetPeriodByID: func(id uuid.UUID) (models.PayrollPeriod, error) {
			return models.PayrollPeriod{}, errors.New("not found")
		},
	}

	err := handler.RunPayroll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "payroll period not found")
}

func TestRunPayroll_ServiceError(t *testing.T) {
	e := echo.New()
	periodID := uuid.New()
	body := map[string]string{"period_id": periodID.String()}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/run-payroll", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := new(mocks.MockPayrollService)
	mockService.On("RunPayroll", models.PayrollPeriod{ID: periodID}).Return(errors.New("processing error"))

	handler := &PayrollHandler{
		Service: mockService,
		GetPeriodByID: func(id uuid.UUID) (models.PayrollPeriod, error) {
			return models.PayrollPeriod{ID: id}, nil
		},
	}

	err := handler.RunPayroll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "processing error")
	mockService.AssertExpectations(t)
}

func TestRunPayroll_InvalidUUID(t *testing.T) {
	e := echo.New()
	body := map[string]string{"period_id": "not-a-uuid"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/run-payroll", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &PayrollHandler{}

	err := handler.RunPayroll(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "invalid period ID")
}

func TestCreatePayrollPeriod_Unauthorized(t *testing.T) {
	e := echo.New()

	body := map[string]string{
		"start_date": "2025-06-01",
		"end_date":   "2025-06-30",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/payroll-period", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &PayrollHandler{Service: nil} // service won't be called

	err := handler.CreatePayrollPeriod(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "unauthorized")
}
