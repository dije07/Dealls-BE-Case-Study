package repositories

import (
	"time"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
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

func (r *OvertimeRepoImpl) CreateOvertime(userID uuid.UUID, hours int, date time.Time) error {
	overtime := models.Overtime{
		ID:     uuid.New(),
		UserID: userID,
		Hours:  hours,
		Date:   date,
	}
	return database.DB.Create(&overtime).Error
}

func (r *OvertimeRepoImpl) GetOvertimeHistory(userID uuid.UUID) ([]models.Overtime, error) {
	var list []models.Overtime
	err := database.DB.Where("user_id = ?", userID).Order("date desc").Find(&list).Error
	return list, err
}
