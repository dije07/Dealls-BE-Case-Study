package repositories

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/glebarez/sqlite" // ✅ pure Go SQLite driver
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestOvertimeRepo(t *testing.T) *OvertimeRepoImpl {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&models.Overtime{})
	if err != nil {
		t.Fatalf("failed to migrate overtime table: %v", err)
	}

	database.DB = db
	return NewOvertimeRepository()
}

func TestOvertime_CreateAndExists(t *testing.T) {
	repo := setupTestOvertimeRepo(t)
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := repo.CreateOvertime(c, userID, 2, today)
	assert.NoError(t, err)

	exists := repo.OvertimeExists(userID, today)
	assert.True(t, exists)
}

func TestOvertime_GetOvertimeHistory(t *testing.T) {
	repo := setupTestOvertimeRepo(t)
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	_ = repo.CreateOvertime(c, userID, 2, today)
	_ = repo.CreateOvertime(c, userID, 3, today.AddDate(0, 0, -1))

	history, err := repo.GetOvertimeHistory(userID)
	assert.NoError(t, err)
	assert.Len(t, history, 2)
}
