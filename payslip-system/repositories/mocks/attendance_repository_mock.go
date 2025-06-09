package mocks

import (
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockAttendanceRepo struct {
	mock.Mock
}

func (m *MockAttendanceRepo) AttendanceExists(userID uuid.UUID, date time.Time) bool {
	args := m.Called(userID, date)
	return args.Bool(0)
}

func (m *MockAttendanceRepo) CreateAttendance(userID uuid.UUID, date time.Time) error {
	args := m.Called(userID, date)
	return args.Error(0)
}

func (m *MockAttendanceRepo) GetAttendanceHistory(userID uuid.UUID) ([]models.Attendance, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Attendance), args.Error(1)
}
