package seeder

import (
	"fmt"

	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
)

func SeedRoles() {
	roles := []models.Role{
		{Name: "admin"},
		{Name: "employee"},
	}

	for _, role := range roles {
		var exists models.Role
		if err := database.DB.Where("name = ?", role.Name).First(&exists).Error; err != nil {
			database.DB.Create(&role)
			fmt.Printf("✅ Role seeded: %s\n", role.Name)
		}
	}
}
