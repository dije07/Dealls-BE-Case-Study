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
		// Ensure request ID exists
		reqID := c.Request().Header.Get(echo.HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
			c.Request().Header.Set(echo.HeaderXRequestID, reqID)
		}
		c.Response().Header().Set(echo.HeaderXRequestID, reqID)
		c.Set("request_id", reqID)

		// Call next
		err := next(c)

		// Only log write actions
		method := c.Request().Method
		if method != echo.POST && method != echo.PUT && method != echo.DELETE {
			return err
		}

		// Extract user_id safely
		var userID uuid.UUID
		switch v := c.Get("user_id").(type) {
		case string:
			parsed, err := uuid.Parse(v)
			if err == nil {
				userID = parsed
			}
		case uuid.UUID:
			userID = v
		}

		entityID, _ := c.Get("entity_id").(uuid.UUID)

		log := models.AuditLog{
			ID:          uuid.New(),
			Action:      method,
			Entity:      c.Path(),
			EntityID:    entityID,
			PerformedBy: userID,
			IP:          c.RealIP(),
			RequestID:   reqID,
			CreatedAt:   time.Now(),
		}

		_ = database.DB.Create(&log)
		return err
	}
}
