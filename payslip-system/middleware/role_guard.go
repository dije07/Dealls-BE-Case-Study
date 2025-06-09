package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RequireRole(expectedRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role := c.Get("role")
			if role != expectedRole {
				return c.JSON(http.StatusForbidden, echo.Map{
					"error": "Access denied: insufficient role",
				})
			}
			return next(c)
		}
	}
}
