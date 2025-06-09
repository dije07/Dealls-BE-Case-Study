package repositories

import (
	"testing"
	"time"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/glebarez/sqlite" // ✅ pure Go SQLite driver
	"github.com/google/uuid"
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

	err := repo.CreateOvertime(userID, 2, today)
	assert.NoError(t, err)

	exists := repo.OvertimeExists(userID, today)
	assert.True(t, exists)
}

func TestOvertime_GetOvertimeHistory(t *testing.T) {
	repo := setupTestOvertimeRepo(t)
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	_ = repo.CreateOvertime(userID, 2, today)
	_ = repo.CreateOvertime(userID, 3, today.AddDate(0, 0, -1))

	history, err := repo.GetOvertimeHistory(userID)
	assert.NoError(t, err)
	assert.Len(t, history, 2)
}
