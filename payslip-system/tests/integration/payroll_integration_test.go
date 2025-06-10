package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreatePayrollPeriodIntegration(t *testing.T) {
	body := `{"start_date":"2025-06-01","end_date":"2025-06-30"}`
	req := httptest.NewRequest(http.MethodPost, "/api/payroll-period", strings.NewReader(body))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+adminToken)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "payroll period created")
}

func TestRunPayrollIntegration(t *testing.T) {
	// Ensure there's an active payroll period
	var period models.PayrollPeriod
	err := database.DB.Order("created_at desc").First(&period).Error
	assert.NoError(t, err)

	// Prepare request
	body := fmt.Sprintf(`{"period_id":"%s"}`, period.ID.String())
	req := httptest.NewRequest(http.MethodPost, "/api/run-payroll", strings.NewReader(body))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+adminToken)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "payroll processed")
}
