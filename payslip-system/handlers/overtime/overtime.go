package handlers

import (
	"net/http"

	"github.com/dije07/payslip-system/config"
	"github.com/dije07/payslip-system/services/interfaces"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OvertimeHandler struct {
	Service interfaces.OvertimeService
}

type OvertimeRequest struct {
	Hours int `json:"hours"`
}

func (h *OvertimeHandler) SubmitOvertime(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	userID, _ := uuid.Parse(userIDStr)

	var req OvertimeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	if err := h.Service.SubmitOvertime(c, userID, req.Hours); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "overtime submitted"})
}

func (h *OvertimeHandler) GetMyOvertime(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	userID, _ := uuid.Parse(userIDStr)

	overtimes, err := h.Service.GetMyOvertime(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch overtime history"})
	}

	type OvertimeResponse struct {
		ID    uuid.UUID `json:"id"`
		Date  string    `json:"date"`
		Hours int       `json:"hours"`
		Time  string    `json:"submitted_at"`
	}

	var result []OvertimeResponse
	for _, ot := range overtimes {
		result = append(result, OvertimeResponse{
			ID:    ot.ID,
			Date:  ot.Date.Format("02 January 2006"),
			Hours: ot.Hours,
			Time:  ot.CreatedAt.Format(config.DateTimeFormat),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{"overtime": result})
}
