package repositories

import (
	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

type PayslipRepoImpl struct{}

func NewPayslipRepository() *PayslipRepoImpl {
	return &PayslipRepoImpl{}
}

func (r *PayslipRepoImpl) GetPayslip(userID, periodID uuid.UUID) (*models.Payslip, error) {
	var p models.Payslip
	err := database.DB.Where("user_id = ? AND period_id = ?", userID, periodID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PayslipRepoImpl) GetPayslipsByPeriod(periodID uuid.UUID) ([]models.Payslip, error) {
	var list []models.Payslip
	err := database.DB.Where("period_id = ?", periodID).Find(&list).Error
	return list, err
}
