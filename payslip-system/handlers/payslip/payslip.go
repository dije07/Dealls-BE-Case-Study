package handlers

import (
	"net/http"

	"github.com/dije07/payslip-system/services/interfaces"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PayslipHandler struct {
	Service interfaces.PayslipService
}

func (h *PayslipHandler) GetMyPayslip(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	userID, _ := uuid.Parse(userIDStr)

	periodIDStr := c.Param("period_id")
	periodID, err := uuid.Parse(periodIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid payroll period ID"})
	}

	payslip, err := h.Service.GetEmployeePayslip(userID, periodID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "payslip not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"attendance_days": payslip.AttendanceDays,
		"base_salary":     payslip.BaseSalary,
		"overtime_hours":  payslip.OvertimeHours,
		"overtime_pay":    payslip.OvertimePay,
		"reimbursement":   payslip.Reimbursement,
		"take_home_pay":   payslip.TakeHomePay,
	})
}

func (h *PayslipHandler) GetPayslipSummary(c echo.Context) error {
	periodIDStr := c.Param("period_id")
	periodID, err := uuid.Parse(periodIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid period ID"})
	}

	payslips, total, err := h.Service.GetPayslipSummary(periodID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch summary"})
	}

	type Summary struct {
		UserID        uuid.UUID `json:"user_id"`
		BaseSalary    float64   `json:"base_salary"`
		Attendance    int       `json:"attendance_days"`
		Overtime      int       `json:"overtime_hours"`
		OvertimePay   float64   `json:"overtime_pay"`
		Reimbursement float64   `json:"reimbursement"`
		TakeHome      float64   `json:"take_home_pay"`
	}

	var result []Summary
	for _, p := range payslips {
		result = append(result, Summary{
			UserID:        p.UserID,
			BaseSalary:    p.BaseSalary,
			Attendance:    p.AttendanceDays,
			Overtime:      p.OvertimeHours,
			OvertimePay:   p.OvertimePay,
			Reimbursement: p.Reimbursement,
			TakeHome:      p.TakeHomePay,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"summary":           result,
		"total_take_home":   total,
		"total_employees":   len(result),
		"payroll_period_id": periodID,
	})
}
