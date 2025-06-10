package handlers

import (
	"net/http"

	"github.com/dije07/payslip-system/config"
	"github.com/dije07/payslip-system/services/interfaces"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AttendanceHandler struct {
	Service interfaces.AttendanceService
}

func (h *AttendanceHandler) SubmitAttendance(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	userID, _ := uuid.Parse(userIDStr)

	if err := h.Service.SubmitAttendance(c, userID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "attendance submitted successfully"})
}

func (h *AttendanceHandler) GetMyAttendance(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	userID, _ := uuid.Parse(userIDStr)

	attendances, err := h.Service.GetMyAttendance(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch attendance"})
	}

	type AttendanceResponse struct {
		ID   uuid.UUID `json:"id"`
		Time string    `json:"time"`
	}

	var result []AttendanceResponse
	for _, a := range attendances {
		result = append(result, AttendanceResponse{
			ID:   a.ID,
			Time: a.CreatedAt.Format(config.DateTimeFormat),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{"attendance": result})
}
