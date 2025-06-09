package seeder

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/dije07/payslip-system/config"
	"github.com/dije07/payslip-system/database"
	"github.com/dije07/payslip-system/models"
)

func SeedUsers() {
	var userCount int64
	database.DB.Model(&models.User{}).Count(&userCount)
	if userCount > 0 {
		fmt.Println("✅ Users already seeded")
		return
	}

	var adminRole, employeeRole models.Role
	database.DB.First(&adminRole, "name = ?", "admin")
	database.DB.First(&employeeRole, "name = ?", "employee")

	// Admin
	pwAdmin, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		ID:       uuid.New(),
		Username: "admin",
		Password: string(pwAdmin),
		RoleID:   adminRole.ID,
	}
	database.DB.Create(&admin)

	// Employees
	for i := 0; i < 100; i++ {
		uname := strings.ToLower(faker.Username()) + fmt.Sprintf("%03d", i)
		pwHash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		emp := models.User{
			ID:       uuid.New(),
			Username: uname,
			Password: string(pwHash),
			RoleID:   employeeRole.ID,
			Salary:   float64(rand.Intn(config.MaxSalary-config.MinSalary) + config.MinSalary),
		}

		database.DB.Create(&emp)
	}

	fmt.Println("✅ Seeded 1 admin and 100 employees with roles")
}
