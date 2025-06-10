package repositories

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestPayrollRepo(t *testing.T) *PayrollRepoImpl {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(
		&models.PayrollPeriod{},
		&models.User{},
		&models.Payslip{},
		&models.Attendance{},
		&models.Overtime{},
		&models.Reimbursement{},
	)
	if err != nil {
		t.Fatalf("failed to migrate test schema: %v", err)
	}

	database.DB = db
	return NewPayrollRepository()
}

func TestPayrollPeriod_CreateAndExists(t *testing.T) {
	repo := setupTestPayrollRepo(t)
	start := time.Now()
	end := start.AddDate(0, 0, 7)

	exists := repo.PayrollPeriodExists(start, end)
	assert.False(t, exists)

	id := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err := repo.CreatePayrollPeriod(c, id, start, end)
	assert.NoError(t, err)

	exists = repo.PayrollPeriodExists(start, end)
	assert.True(t, exists)
}

func TestGetAllEmployees(t *testing.T) {
	repo := setupTestPayrollRepo(t)

	// Insert admin (role_id 1) and employee (role_id 2)
	database.DB.Create(&models.User{ID: uuid.New(), Username: "admin", RoleID: 1, Salary: 9000000})
	database.DB.Create(&models.User{ID: uuid.New(), Username: "employee", RoleID: 2, Salary: 10000000})

	employees, err := repo.GetAllEmployees()
	assert.NoError(t, err)
	assert.Len(t, employees, 1)
	assert.Equal(t, 2, int(employees[0].RoleID))
}

func TestSaveAndClosePayslip(t *testing.T) {
	repo := setupTestPayrollRepo(t)
	period := models.PayrollPeriod{ID: uuid.New(), StartDate: time.Now(), EndDate: time.Now().AddDate(0, 0, 5)}
	userID := uuid.New()

	_ = database.DB.Create(&period)

	p := models.Payslip{
		ID:             uuid.New(),
		UserID:         userID,
		PeriodID:       period.ID,
		BaseSalary:     5000000,
		AttendanceDays: 5,
		OvertimeHours:  2,
		Reimbursement:  250000,
		TakeHomePay:    6000000,
	}

	err := repo.SavePayslip(p)
	assert.NoError(t, err)

	err = repo.ClosePayrollPeriod(period.ID)
	assert.NoError(t, err)

	var updated models.PayrollPeriod
	database.DB.First(&updated, "id = ?", period.ID)
	assert.True(t, updated.IsClosed)
}

func TestCountAttendances(t *testing.T) {
	repo := setupTestPayrollRepo(t)
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	// Seed attendances
	database.DB.Create(&models.Attendance{ID: uuid.New(), UserID: userID, Date: today})
	database.DB.Create(&models.Attendance{ID: uuid.New(), UserID: userID, Date: today.AddDate(0, 0, -1)})

	count, err := repo.CountAttendances(userID, today.AddDate(0, 0, -2), today)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestSumOvertimeHours(t *testing.T) {
	repo := setupTestPayrollRepo(t)
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	database.DB.Create(&models.Overtime{ID: uuid.New(), UserID: userID, Hours: 2, Date: today})
	database.DB.Create(&models.Overtime{ID: uuid.New(), UserID: userID, Hours: 1, Date: today.AddDate(0, 0, -1)})

	sum, err := repo.SumOvertimeHours(userID, today.AddDate(0, 0, -2), today)
	assert.NoError(t, err)
	assert.Equal(t, 3, sum)
}

func TestSumReimbursements(t *testing.T) {
	repo := setupTestPayrollRepo(t)
	userID := uuid.New()
	now := time.Now()

	database.DB.Create(&models.Reimbursement{ID: uuid.New(), UserID: userID, Amount: 50000, CreatedAt: now})
	database.DB.Create(&models.Reimbursement{ID: uuid.New(), UserID: userID, Amount: 100000, CreatedAt: now.Add(-time.Hour)})

	total, err := repo.SumReimbursements(userID, now.Add(-24*time.Hour), now.Add(1*time.Hour))
	assert.NoError(t, err)
	assert.Equal(t, 150000.0, total)
}
