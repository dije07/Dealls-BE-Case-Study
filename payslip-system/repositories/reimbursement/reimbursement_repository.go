package repositories

import (
	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

type ReimbursementRepoImpl struct{}

func NewReimbursementRepository() *ReimbursementRepoImpl {
	return &ReimbursementRepoImpl{}
}

func (r *ReimbursementRepoImpl) CreateReimbursement(userID uuid.UUID, amount float64, description string) error {
	rb := models.Reimbursement{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      amount,
		Description: description,
	}
	return database.DB.Create(&rb).Error
}

func (r *ReimbursementRepoImpl) GetReimbursementsByUser(userID uuid.UUID) ([]models.Reimbursement, error) {
	var list []models.Reimbursement
	err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error
	return list, err
}
