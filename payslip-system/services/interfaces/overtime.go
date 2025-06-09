package interfaces

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
)

type OvertimeService interface {
	SubmitOvertime(userID uuid.UUID, hours int) error
	GetMyOvertime(uuid.UUID) ([]models.Overtime, error)
}
