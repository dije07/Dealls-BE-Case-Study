package services

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/repositories/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSubmitAttendance_Success(t *testing.T) {
	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")

	mockRepo := new(mocks.MockAttendanceRepo)
	userID := uuid.New()
	today := time.Now().Truncate(24 * time.Hour)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	mockRepo.On("AttendanceExists", userID, today).Return(false)
	mockRepo.On("CreateAttendance", c, userID, today).Return(nil)

	service := &AttendanceServiceImpl{Repo: mockRepo}
	err := service.SubmitAttendance(c, userID)

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

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	service := &AttendanceServiceImpl{Repo: mockRepo}
	err := service.SubmitAttendance(c, userID)

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
	// ⏰ Force a Saturday
	nowFunc = func() time.Time {
		return time.Date(2025, 6, 7, 9, 0, 0, 0, time.UTC)
	}
	defer func() { nowFunc = time.Now }()

	os.Unsetenv("TEST_MODE")
	userID := uuid.New()

	mockRepo := new(mocks.MockAttendanceRepo)
	service := &AttendanceServiceImpl{Repo: mockRepo}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := service.SubmitAttendance(c, userID)

	assert.EqualError(t, err, "cannot submit attendance on weekends")
	mockRepo.AssertNotCalled(t, "AttendanceExists")
	mockRepo.AssertNotCalled(t, "CreateAttendance")
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
