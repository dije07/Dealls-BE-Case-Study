package main

import (
	"log"

	"github.com/dije07/payslip-system/database"
	custommiddleware "github.com/dije07/payslip-system/middleware"
	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/routes"
	"github.com/dije07/payslip-system/seeder"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	database.Connect()
	err := database.DB.AutoMigrate(&models.AuditLog{}, &models.Payslip{}, &models.PayrollPeriod{}, &models.Reimbursement{}, &models.Overtime{}, &models.Attendance{}, &models.Role{}, &models.User{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	seeder.SeedRoles()
	seeder.SeedUsers()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(custommiddleware.AuditLogger)

	routes.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
