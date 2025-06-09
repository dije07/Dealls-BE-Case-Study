package repositories

import (
	"testing"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestPayslipRepo(t *testing.T) *PayslipRepoImpl {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&models.Payslip{})
	if err != nil {
		t.Fatalf("failed to migrate payslip table: %v", err)
	}

	database.DB = db
	return NewPayslipRepository()
}

func TestGetPayslip(t *testing.T) {
	repo := setupTestPayslipRepo(t)
	userID := uuid.New()
	periodID := uuid.New()

	expected := models.Payslip{
		ID:             uuid.New(),
		UserID:         userID,
		PeriodID:       periodID,
		TakeHomePay:    8500000,
		AttendanceDays: 20,
	}

	err := database.DB.Create(&expected).Error
	assert.NoError(t, err)

	p, err := repo.GetPayslip(userID, periodID)
	assert.NoError(t, err)
	assert.Equal(t, expected.TakeHomePay, p.TakeHomePay)
	assert.Equal(t, expected.AttendanceDays, p.AttendanceDays)
}

func TestGetPayslipsByPeriod(t *testing.T) {
	repo := setupTestPayslipRepo(t)
	periodID := uuid.New()

	payslips := []models.Payslip{
		{ID: uuid.New(), UserID: uuid.New(), PeriodID: periodID, TakeHomePay: 7000000},
		{ID: uuid.New(), UserID: uuid.New(), PeriodID: periodID, TakeHomePay: 6500000},
	}

	for _, p := range payslips {
		err := database.DB.Create(&p).Error
		assert.NoError(t, err)
	}

	result, err := repo.GetPayslipsByPeriod(periodID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 7000000.0, result[0].TakeHomePay)
}

func TestGetPayslip_NotFound(t *testing.T) {
	repo := setupTestPayslipRepo(t)
	userID := uuid.New()
	periodID := uuid.New()

	p, err := repo.GetPayslip(userID, periodID)

	assert.Nil(t, p)
	assert.Error(t, err)
}
