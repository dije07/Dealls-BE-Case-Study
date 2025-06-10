package repositories

import (
	"time"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OvertimeRepoImpl struct{}

func NewOvertimeRepository() *OvertimeRepoImpl {
	return &OvertimeRepoImpl{}
}

func (r *OvertimeRepoImpl) OvertimeExists(userID uuid.UUID, date time.Time) bool {
	var ot models.Overtime
	err := database.DB.Where("user_id = ? AND date = ?", userID, date).First(&ot).Error
	return err == nil
}

func (r *OvertimeRepoImpl) CreateOvertime(c echo.Context, userID uuid.UUID, hours int, date time.Time) error {
	overtime := models.Overtime{
		ID:        uuid.New(),
		UserID:    userID,
		Hours:     hours,
		Date:      date,
		CreatedBy: userID,
		UpdatedBy: userID,
		IPAddress: c.RealIP(),
	}
	c.Set("entity_id", overtime.ID)
	return database.DB.Create(&overtime).Error
}

func (r *OvertimeRepoImpl) GetOvertimeHistory(userID uuid.UUID) ([]models.Overtime, error) {
	var list []models.Overtime
	err := database.DB.Where("user_id = ?", userID).Order("date desc").Find(&list).Error
	return list, err
}
