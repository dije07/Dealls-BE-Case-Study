package repositories

import (
	"time"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

type AttendanceRepoImpl struct{}

func NewAttendanceRepository() *AttendanceRepoImpl {
	return &AttendanceRepoImpl{}
}

func (r *AttendanceRepoImpl) AttendanceExists(userID uuid.UUID, date time.Time) bool {
	var a models.Attendance
	err := database.DB.Where("user_id = ? AND date = ?", userID, date).First(&a).Error
	return err == nil
}

func (r *AttendanceRepoImpl) CreateAttendance(userID uuid.UUID, date time.Time) error {
	a := models.Attendance{
		ID:     uuid.New(), // ✅ ensure unique ID
		UserID: userID,
		Date:   date,
	}
	return database.DB.Create(&a).Error
}

func (r *AttendanceRepoImpl) GetAttendanceHistory(userID uuid.UUID) ([]models.Attendance, error) {
	var list []models.Attendance
	err := database.DB.Where("user_id = ?", userID).Order("date desc").Find(&list).Error
	return list, err
}
