package handlers

import (
	"net/http"
	"time"

	"github.com/dije07/payslip-system/config"
	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/services/interfaces"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PayrollHandler struct {
	Service       interfaces.PayrollService
	GetPeriodByID func(uuid.UUID) (models.PayrollPeriod, error) // injectable for mocking
}

type PayrollPeriodRequest struct {
	StartDate string `json:"start_date"` // e.g. "2025-06-01"
	EndDate   string `json:"end_date"`   // e.g. "2025-06-30"
}

func (h *PayrollHandler) CreatePayrollPeriod(c echo.Context) error {
	var req PayrollPeriodRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	layout := config.DateFormat
	start, err1 := time.Parse(layout, req.StartDate)
	end, err2 := time.Parse(layout, req.EndDate)

	if err1 != nil || err2 != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid date format (use YYYY-MM-DD)"})
	}

	if err := h.Service.CreatePayrollPeriod(start, end); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "payroll period created"})
}

func (h *PayrollHandler) RunPayroll(c echo.Context) error {
	type Req struct {
		PeriodID string `json:"period_id"`
	}

	var req Req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	periodID, err := uuid.Parse(req.PeriodID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid period ID"})
	}

	period, err := h.GetPeriodByID(periodID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "payroll period not found"})
	}

	if err := h.Service.RunPayroll(period); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "payroll processed successfully"})
}
