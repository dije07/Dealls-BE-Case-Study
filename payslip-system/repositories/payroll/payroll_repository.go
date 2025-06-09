package repositories

import (
	"time"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
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

func (r *PayrollRepoImpl) CreatePayrollPeriod(start, end time.Time) error {
	p := models.PayrollPeriod{
		ID:        uuid.New(),
		StartDate: start,
		EndDate:   end,
		IsClosed:  false,
	}
	return database.DB.Create(&p).Error
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
