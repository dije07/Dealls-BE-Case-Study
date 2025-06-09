package services

import (
	"errors"
	"os"
	"time"

	"github.com/dije07/payslip-system/models"
	repositoriesInterface "github.com/dije07/payslip-system/repositories/interfaces"
	"github.com/dije07/payslip-system/services/interfaces"
	"github.com/google/uuid"
)

type AttendanceServiceImpl struct {
	Repo repositoriesInterface.AttendanceRepository // good: interface
}

func NewAttendanceService(repo repositoriesInterface.AttendanceRepository) interfaces.AttendanceService {
	return &AttendanceServiceImpl{Repo: repo}
}

func (s *AttendanceServiceImpl) SubmitAttendance(userID uuid.UUID) error {
	today := time.Now().Truncate(24 * time.Hour)
	if today.Weekday() == time.Saturday || today.Weekday() == time.Sunday {
		if os.Getenv("TEST_MODE") != "true" {
			return errors.New("cannot submit attendance on weekends")
		}
	}
	if s.Repo.AttendanceExists(userID, today) {
		return errors.New("attendance already submitted for today")
	}
	return s.Repo.CreateAttendance(userID, today)
}

func (s *AttendanceServiceImpl) GetMyAttendance(userID uuid.UUID) ([]models.Attendance, error) {
	return s.Repo.GetAttendanceHistory(userID)
}
