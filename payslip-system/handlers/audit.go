package handlers

import (
	"net/http"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/labstack/echo/v4"
)

func GetAuditLogs(c echo.Context) error {
	var logs []models.AuditLog
	if err := database.DB.Order("created_at desc").Limit(100).Find(&logs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch logs"})
	}
	return c.JSON(http.StatusOK, logs)
}
