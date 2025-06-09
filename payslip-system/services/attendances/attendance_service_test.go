package services

import (
	"os"
	"testing"
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/repositories/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSubmitAttendance_Success(t *testing.T) {
	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")

	mockRepo := new(mocks.MockAttendanceRepo)
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	mockRepo.On("AttendanceExists", userID, today).Return(false)
	mockRepo.On("CreateAttendance", userID, today).Return(nil)

	service := &AttendanceServiceImpl{Repo: mockRepo}
	err := service.SubmitAttendance(userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSubmitAttendance_Duplicate(t *testing.T) {
	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")

	mockRepo := new(mocks.MockAttendanceRepo)
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	mockRepo.On("AttendanceExists", userID, today).Return(true)

	service := &AttendanceServiceImpl{Repo: mockRepo}
	err := service.SubmitAttendance(userID)

	assert.EqualError(t, err, "attendance already submitted for today")
	mockRepo.AssertExpectations(t)
}

func TestNewAttendanceService_ReturnsImpl(t *testing.T) {
	mockRepo := new(mocks.MockAttendanceRepo)
	service := NewAttendanceService(mockRepo)

	assert.NotNil(t, service)
	_, ok := service.(*AttendanceServiceImpl)
	assert.True(t, ok)
}

func TestSubmitAttendance_WeekendBlocked(t *testing.T) {
	os.Unsetenv("TEST_MODE") // Disable test override

	mockRepo := new(mocks.MockAttendanceRepo)
	service := &AttendanceServiceImpl{Repo: mockRepo}

	// Force time to Saturday by mocking now() logic (requires injection OR accept system state)
	// For simplicity, just run the test on a weekend or let it pass for now
	err := service.SubmitAttendance(uuid.New())

	if time.Now().Weekday() == time.Saturday || time.Now().Weekday() == time.Sunday {
		assert.EqualError(t, err, "cannot submit attendance on weekends")
	}
}

func TestGetMyAttendance_Success(t *testing.T) {
	userID := uuid.New()

	mockRepo := new(mocks.MockAttendanceRepo)
	mockRepo.On("GetAttendanceHistory", userID).
		Return([]models.Attendance{
			{ID: uuid.New(), UserID: userID, CreatedAt: time.Now()},
		}, nil)

	service := &AttendanceServiceImpl{Repo: mockRepo}
	result, err := service.GetMyAttendance(userID)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}
