package services

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/repositories/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePayrollPeriod_Success(t *testing.T) {
	start := time.Now()
	end := start.AddDate(0, 0, 5)
	id := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	mockRepo := new(mocks.MockPayrollRepo)
	mockRepo.On("PayrollPeriodExists", start, end).Return(false)
	mockRepo.On("CreatePayrollPeriod", c, id, start, end).Return(nil)

	service := NewPayrollService(mockRepo)
	err := service.CreatePayrollPeriod(c, id, start, end)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreatePayrollPeriod_EndBeforeStart(t *testing.T) {
	start := time.Now()
	end := start.AddDate(0, 0, -1)
	id := uuid.New()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	service := NewPayrollService(new(mocks.MockPayrollRepo))
	err := service.CreatePayrollPeriod(c, id, start, end)

	assert.EqualError(t, err, "end date cannot be before start date")
}

func TestCreatePayrollPeriod_Duplicate(t *testing.T) {
	start := time.Now()
	end := start.AddDate(0, 0, 7)
	id := uuid.New()

	mockRepo := new(mocks.MockPayrollRepo)
	mockRepo.On("PayrollPeriodExists", start, end).Return(true)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	service := NewPayrollService(mockRepo)
	err := service.CreatePayrollPeriod(c, id, start, end)

	assert.EqualError(t, err, "payroll period already exists")
	mockRepo.AssertExpectations(t)
}

func TestRunPayroll_AlreadyClosed(t *testing.T) {
	service := NewPayrollService(new(mocks.MockPayrollRepo))
	period := models.PayrollPeriod{ID: uuid.New(), IsClosed: true}

	err := service.RunPayroll(period)
	assert.EqualError(t, err, "payroll already processed for this period")
}

func TestRunPayroll_Success(t *testing.T) {
	period := models.PayrollPeriod{
		ID:        uuid.New(),
		IsClosed:  false,
		StartDate: time.Now().AddDate(0, 0, -5),
		EndDate:   time.Now(),
	}

	emp := models.User{ID: uuid.New(), Salary: 10000000}
	mockRepo := new(mocks.MockPayrollRepo)

	mockRepo.On("GetAllEmployees").Return([]models.User{emp}, nil)
	mockRepo.On("CountAttendances", emp.ID, period.StartDate, period.EndDate).Return(5, nil)
	mockRepo.On("SumOvertimeHours", emp.ID, period.StartDate, period.EndDate).Return(3, nil)
	mockRepo.On("SumReimbursements", emp.ID, period.StartDate, period.EndDate).Return(250000.0, nil)
	mockRepo.On("SavePayslip", mock.Anything).Return(nil)
	mockRepo.On("ClosePayrollPeriod", period.ID).Return(nil)

	service := NewPayrollService(mockRepo)
	err := service.RunPayroll(period)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestNewPayrollService_ReturnsImplementation(t *testing.T) {
	mockRepo := new(mocks.MockPayrollRepo)
	service := NewPayrollService(mockRepo)

	assert.NotNil(t, service)
	_, ok := service.(*PayrollServiceImpl)
	assert.True(t, ok)
}

func TestRunPayroll_FailsOnGetAllEmployees(t *testing.T) {
	period := models.PayrollPeriod{
		ID:        uuid.New(),
		IsClosed:  false,
		StartDate: time.Now().AddDate(0, 0, -5),
		EndDate:   time.Now(),
	}

	mockRepo := new(mocks.MockPayrollRepo)
	mockRepo.On("GetAllEmployees").Return([]models.User{}, errors.New("db error"))

	service := NewPayrollService(mockRepo)
	err := service.RunPayroll(period)

	assert.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}
