package integration

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSubmitAttendanceIntegration(t *testing.T) {
	// Temporarily set test mode to bypass weekend rule
	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")

	req := httptest.NewRequest(http.MethodPost, "/api/attendance", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+employeeToken)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	// Serve the request
	e.ServeHTTP(rec, req)

	// Check
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "attendance submitted")
}
