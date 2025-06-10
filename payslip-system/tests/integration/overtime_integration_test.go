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

func TestSubmitOvertimeIntegration(t *testing.T) {
	// Set test mode for consistent behavior
	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")

	body := map[string]interface{}{
		"hours": 2,
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/overtime", bytes.NewBuffer(bodyBytes))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+employeeToken)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "overtime submitted")
}
