package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSubmitReimbursementIntegration(t *testing.T) {
	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")

	body := map[string]interface{}{
		"amount":      150000,
		"description": "Transport reimbursement",
	}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/reimbursement", bytes.NewReader(payload))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+employeeToken)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "reimbursement submitted")
}

func TestGetMyReimbursementsIntegration(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/reimbursement", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+employeeToken)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "reimbursements")
}
