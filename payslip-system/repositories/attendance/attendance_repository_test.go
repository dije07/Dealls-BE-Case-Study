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

func setupTestAttendanceRepo(t *testing.T) *AttendanceRepoImpl {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&models.Attendance{})
	if err != nil {
		t.Fatalf("failed to migrate attendance table: %v", err)
	}

	database.DB = db
	return NewAttendanceRepository()
}

func TestAttendance_CreateAndExists(t *testing.T) {
	repo := setupTestAttendanceRepo(t)
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	// Act: Create attendance
	err := repo.CreateAttendance(userID, today)
	assert.NoError(t, err)

	// Assert: Should exist
	exists := repo.AttendanceExists(userID, today)
	assert.True(t, exists)
}

func TestAttendance_GetAttendanceHistory(t *testing.T) {
	repo := setupTestAttendanceRepo(t)
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	_ = repo.CreateAttendance(userID, today)
	_ = repo.CreateAttendance(userID, today.AddDate(0, 0, -1))

	history, err := repo.GetAttendanceHistory(userID)

	assert.NoError(t, err)
	assert.Len(t, history, 2)
}
