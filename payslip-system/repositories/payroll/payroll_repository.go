package repositories

import (
	"time"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PayrollRepoImpl struct{}

func NewPayrollRepository() *PayrollRepoImpl {
	return &PayrollRepoImpl{}
}

func (r *PayrollRepoImpl) PayrollPeriodExists(start, end time.Time) bool {
	var p models.PayrollPeriod
	err := database.DB.Where("start_date = ? AND end_date = ?", start, end).First(&p).Error
	return err == nil
}

func (r *PayrollRepoImpl) CreatePayrollPeriod(c echo.Context, userID uuid.UUID, start, end time.Time) error {
	payroll := models.PayrollPeriod{
		ID:        uuid.New(),
		StartDate: start,
		EndDate:   end,
		IsClosed:  false,
		CreatedBy: userID,
		UpdatedBy: userID,
		IPAddress: c.RealIP(),
	}
	c.Set("entity_id", payroll.ID)
	return database.DB.Create(&payroll).Error
}

func (r *PayrollRepoImpl) GetAllEmployees() ([]models.User, error) {
	var users []models.User
	err := database.DB.Where("role_id = ?", 2).Find(&users).Error
	return users, err
}

func (r *PayrollRepoImpl) CountAttendances(userID uuid.UUID, start, end time.Time) (int, error) {
	var count int64
	err := database.DB.Model(&models.Attendance{}).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, start, end).
		Count(&count).Error
	return int(count), err
}

func (r *PayrollRepoImpl) SumOvertimeHours(userID uuid.UUID, start, end time.Time) (int, error) {
	var total int64
	err := database.DB.Model(&models.Overtime{}).
		Where("user_id = ? AND date BETWEEN ? AND ?", userID, start, end).
		Select("SUM(hours)").Scan(&total).Error
	return int(total), err
}

func (r *PayrollRepoImpl) SumReimbursements(userID uuid.UUID, start, end time.Time) (float64, error) {
	var total float64
	err := database.DB.Model(&models.Reimbursement{}).
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, start, end).
		Select("SUM(amount)").Scan(&total).Error
	return total, err
}

func (r *PayrollRepoImpl) SavePayslip(p models.Payslip) error {
	return database.DB.Create(&p).Error
}

func (r *PayrollRepoImpl) ClosePayrollPeriod(periodID uuid.UUID) error {
	return database.DB.Model(&models.PayrollPeriod{}).
		Where("id = ?", periodID).
		Update("is_closed", true).Error
}

func (r *PayrollRepoImpl) GetPeriodByID(id uuid.UUID) (models.PayrollPeriod, error) {
	var period models.PayrollPeriod
	err := database.DB.First(&period, "id = ?", id).Error
	return period, err
}
