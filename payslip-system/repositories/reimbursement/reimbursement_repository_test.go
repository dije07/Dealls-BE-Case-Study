package repositories

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestReimbursementRepo(t *testing.T) *ReimbursementRepoImpl {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&models.Reimbursement{})
	if err != nil {
		t.Fatalf("failed to migrate reimbursement table: %v", err)
	}

	database.DB = db
	return NewReimbursementRepository()
}

func TestCreateReimbursement(t *testing.T) {
	repo := setupTestReimbursementRepo(t)
	userID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := repo.CreateReimbursement(c, userID, 200000, "Transport")
	assert.NoError(t, err)

	var result models.Reimbursement
	err = database.DB.First(&result, "user_id = ?", userID).Error
	assert.NoError(t, err)
	assert.Equal(t, 200000.0, result.Amount)
	assert.Equal(t, "Transport", result.Description)
}

func TestGetReimbursementsByUser(t *testing.T) {
	repo := setupTestReimbursementRepo(t)
	userID := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	repo.CreateReimbursement(c, userID, 100000, "Meals")
	repo.CreateReimbursement(c, userID, 50000, "Taxi")

	list, err := repo.GetReimbursementsByUser(userID)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.ElementsMatch(t, []string{"Meals", "Taxi"}, []string{
		list[0].Description,
		list[1].Description,
	})
}
