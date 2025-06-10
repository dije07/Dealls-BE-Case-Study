package interfaces

import (
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OvertimeService interface {
	SubmitOvertime(c echo.Context, userID uuid.UUID, hours int) error
	GetMyOvertime(uuid.UUID) ([]models.Overtime, error)
}
