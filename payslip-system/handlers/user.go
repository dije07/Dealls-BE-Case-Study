package handlers

import (
	"net/http"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GetMyProfile(c echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	var user models.User
	if err := database.DB.Preload("Role").First(&user, "id = ?", userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role.Name,
		"created":  user.CreatedAt,
	})
}
