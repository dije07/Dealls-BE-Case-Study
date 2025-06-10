package services

import (
	"errors"
	"time"

	"github.com/dije07/payslip-system/models"
	repositoryInterfaces "github.com/dije07/payslip-system/repositories/interfaces"
	"github.com/dije07/payslip-system/services/interfaces"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OvertimeServiceImpl struct {
	Repo repositoryInterfaces.OvertimeRepository
}

func NewOvertimeService(repo repositoryInterfaces.OvertimeRepository) interfaces.OvertimeService {
	return &OvertimeServiceImpl{Repo: repo}
}

func (s *OvertimeServiceImpl) SubmitOvertime(c echo.Context, userID uuid.UUID, hours int) error {
	if hours < 1 || hours > 3 {
		return errors.New("overtime must be between 1–3 hours")
	}

	today := time.Now().Truncate(24 * time.Hour)
	if s.Repo.OvertimeExists(userID, today) {
		return errors.New("overtime already submitted for today")
	}

	return s.Repo.CreateOvertime(c, userID, hours, today)
}

func (s *OvertimeServiceImpl) GetMyOvertime(userID uuid.UUID) ([]models.Overtime, error) {
	return s.Repo.GetOvertimeHistory(userID)
}
