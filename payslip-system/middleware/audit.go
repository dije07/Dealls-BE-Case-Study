package middleware

import (
	"time"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func AuditLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Let the request proceed
		err := next(c)

		// Only log POST, PUT, DELETE (not GETs)
		if c.Request().Method != echo.POST && c.Request().Method != echo.PUT && c.Request().Method != echo.DELETE {
			return err
		}

		// Try to get user_id from context (set by JWT middleware)
		userIDVal := c.Get("user_id")
		var userID uuid.UUID
		if userIDStr, ok := userIDVal.(string); ok {
			parsed, err := uuid.Parse(userIDStr)
			if err == nil {
				userID = parsed
			}
		}

		// Save audit log
		log := models.AuditLog{
			Action:      c.Request().Method, // e.g., "POST"
			Entity:      c.Path(),           // e.g., "/api/attendance"
			EntityID:    uuid.Nil,           // optional: can be updated later if needed
			PerformedBy: userID,
			IP:          c.RealIP(),
			RequestID:   c.Response().Header().Get(echo.HeaderXRequestID),
			CreatedAt:   time.Now(),
		}

		_ = database.DB.Create(&log)

		return err
	}
}
