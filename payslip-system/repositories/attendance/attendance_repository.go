package repositories

import (
	"time"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

func (r *AttendanceRepoImpl) CreateAttendance(c echo.Context, userID uuid.UUID, date time.Time) error {
	attendance := models.Attendance{
		ID:        uuid.New(), // ✅ ensure unique ID
		UserID:    userID,
		Date:      date,
		CreatedBy: userID,
		UpdatedBy: userID,
		IPAddress: c.RealIP(),
	}
	c.Set("entity_id", attendance.ID)
	return database.DB.Create(&attendance).Error
}

func (r *AttendanceRepoImpl) GetAttendanceHistory(userID uuid.UUID) ([]models.Attendance, error) {
	var list []models.Attendance
	err := database.DB.Where("user_id = ?", userID).Order("date desc").Find(&list).Error
	return list, err
}
