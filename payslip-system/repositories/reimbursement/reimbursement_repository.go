package repositories

import (
	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReimbursementRepoImpl struct{}

func NewReimbursementRepository() *ReimbursementRepoImpl {
	return &ReimbursementRepoImpl{}
}

func (r *ReimbursementRepoImpl) CreateReimbursement(c echo.Context, userID uuid.UUID, amount float64, description string) error {
	reimbursement := models.Reimbursement{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      amount,
		Description: description,
		CreatedBy:   userID,
		UpdatedBy:   userID,
		IPAddress:   c.RealIP(),
	}
	c.Set("entity_id", reimbursement.ID)
	return database.DB.Create(&reimbursement).Error
}

func (r *ReimbursementRepoImpl) GetReimbursementsByUser(userID uuid.UUID) ([]models.Reimbursement, error) {
	var list []models.Reimbursement
	err := database.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&list).Error
	return list, err
}
