package integration

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/dije07/payslip-system/database"
	authHandler "github.com/dije07/payslip-system/handlers"
	"github.com/dije07/payslip-system/models"
	"github.com/dije07/payslip-system/routes"
	"github.com/dije07/payslip-system/seeder"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

var e *echo.Echo
var adminToken string
var employeeToken string

func init() {
	// Load test env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("⚠️ .env.test not found, using fallback env")
	}

	// Connect to test DB
	database.Connect()

	// Auto-migrate all required models
	err := database.DB.AutoMigrate(
		&models.AuditLog{},
		&models.Payslip{},
		&models.PayrollPeriod{},
		&models.Reimbursement{},
		&models.Overtime{},
		&models.Attendance{},
		&models.Role{},
		&models.User{},
		&models.AuditLog{},
	)
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	if !HasSeeded() {
		seeder.SeedRoles()
		seeder.SeedUsers()
	}

	// Init Echo
	e = echo.New()
	routes.RegisterRoutes(e)

	// Login as admin
	adminToken = login("admin", "admin123")

	// Find first employee username from DB
	var username string
	err = database.DB.Raw("SELECT username FROM users WHERE role_id = 2 LIMIT 1").Scan(&username).Error
	if err != nil || username == "" {
		log.Fatal("Failed to find employee username for login")
	}
	employeeToken = login(username, "password123")
}

func login(username, password string) string {
	body := map[string]string{
		"username": username,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Trigger login handler
	_ = authHandler.Login(c)

	var res map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &res)

	return res["token"]
}

func authRequest(method, path, token string, body any) (*httptest.ResponseRecorder, error) {
	var reqBody *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewReader(b)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	if token != "" {
		req.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	}

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	return rec, nil
}

func HasSeeded() bool {
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	return count > 0
}
