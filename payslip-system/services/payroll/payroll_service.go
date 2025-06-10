package services

import (
	"errors"
	"time"

	"github.com/dije07/payslip-system/config"
	"github.com/dije07/payslip-system/models"
	repositoryInterfaces "github.com/dije07/payslip-system/repositories/interfaces"
	"github.com/dije07/payslip-system/services/interfaces"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PayrollServiceImpl struct {
	Repo repositoryInterfaces.PayrollRepository
}

func NewPayrollService(repo repositoryInterfaces.PayrollRepository) interfaces.PayrollService {
	return &PayrollServiceImpl{Repo: repo}
}

func (s *PayrollServiceImpl) CreatePayrollPeriod(c echo.Context, userID uuid.UUID, start, end time.Time) error {
	if end.Before(start) {
		return errors.New("end date cannot be before start date")
	}
	if s.Repo.PayrollPeriodExists(start, end) {
		return errors.New("payroll period already exists")
	}
	return s.Repo.CreatePayrollPeriod(c, userID, start, end)
}

func (s *PayrollServiceImpl) RunPayroll(period models.PayrollPeriod) error {
	if period.IsClosed {
		return errors.New("payroll already processed for this period")
	}

	days := s.countWeekdays(period.StartDate, period.EndDate)
	employees, err := s.Repo.GetAllEmployees()
	if err != nil {
		return err
	}

	for _, emp := range employees {
		attended, _ := s.Repo.CountAttendances(emp.ID, period.StartDate, period.EndDate)
		overtimeHours, _ := s.Repo.SumOvertimeHours(emp.ID, period.StartDate, period.EndDate)
		reimb, _ := s.Repo.SumReimbursements(emp.ID, period.StartDate, period.EndDate)

		dailyRate := emp.Salary / float64(days)
		base := dailyRate * float64(attended)
		hourlyRate := dailyRate / float64(config.DefaultWorkingHours)
		overtimePay := hourlyRate * 2 * float64(overtimeHours)
		takeHome := base + overtimePay + reimb

		p := models.Payslip{
			ID:             uuid.New(),
			UserID:         emp.ID,
			PeriodID:       period.ID,
			BaseSalary:     base,
			AttendanceDays: attended,
			OvertimeHours:  overtimeHours,
			OvertimePay:    overtimePay,
			Reimbursement:  reimb,
			TakeHomePay:    takeHome,
		}
		_ = s.Repo.SavePayslip(p)
	}

	return s.Repo.ClosePayrollPeriod(period.ID)
}

func (s *PayrollServiceImpl) countWeekdays(start, end time.Time) int {
	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			count++
		}
	}
	return count
}
