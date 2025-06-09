package handlers

import (
	"net/http"

	"github.com/dije07/payslip-system/config"
	"github.com/dije07/payslip-system/services/interfaces"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ReimbursementHandler struct {
	Service interfaces.ReimbursementService
}

type ReimbursementRequest struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

func (h *ReimbursementHandler) SubmitReimbursement(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	userID, _ := uuid.Parse(userIDStr)

	var req ReimbursementRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	if err := h.Service.SubmitReimbursement(userID, req.Amount, req.Description); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "reimbursement submitted"})
}

func (h *ReimbursementHandler) GetMyReimbursements(c echo.Context) error {
	userIDStr, ok := c.Get("user_id").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	userID, _ := uuid.Parse(userIDStr)

	reimbursements, err := h.Service.GetMyReimbursements(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch reimbursements"})
	}

	type Response struct {
		ID          uuid.UUID `json:"id"`
		Amount      float64   `json:"amount"`
		Description string    `json:"description"`
		Time        string    `json:"submitted_at"`
	}

	var result []Response
	for _, r := range reimbursements {
		result = append(result, Response{
			ID:          r.ID,
			Amount:      r.Amount,
			Description: r.Description,
			Time:        r.CreatedAt.Format(config.DateTimeFormat),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{"reimbursements": result})
}
