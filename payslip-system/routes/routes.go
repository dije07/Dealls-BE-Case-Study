package routes

import (
	"github.com/dije07/payslip-system/handlers"
	attendanceHandlerPkg "github.com/dije07/payslip-system/handlers/attendance"
	overtimeHandlerPkg "github.com/dije07/payslip-system/handlers/overtime"
	payrollHandlerPkg "github.com/dije07/payslip-system/handlers/payroll"
	payslipHandlerPkg "github.com/dije07/payslip-system/handlers/payslip"
	reimbursementHandlerPkg "github.com/dije07/payslip-system/handlers/reimbursement"
	"github.com/dije07/payslip-system/middleware"
	attendanceRepositoryPkg "github.com/dije07/payslip-system/repositories/attendance"
	overtimeRepositoryPkg "github.com/dije07/payslip-system/repositories/overtime"
	payrollRepositoryPkg "github.com/dije07/payslip-system/repositories/payroll"
	payslipRepositoryPkg "github.com/dije07/payslip-system/repositories/payslip"
	reimbursementRepositoryPkg "github.com/dije07/payslip-system/repositories/reimbursement"
	attendanceServicePkg "github.com/dije07/payslip-system/services/attendances"
	overtimeServicePkg "github.com/dije07/payslip-system/services/overtime"
	payrollServicePkg "github.com/dije07/payslip-system/services/payroll"
	payslipServicePkg "github.com/dije07/payslip-system/services/payslip"
	reimbursementServicePkg "github.com/dije07/payslip-system/services/reimbursement"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	// Repositories
	attendanceRepo := attendanceRepositoryPkg.NewAttendanceRepository()
	overtimeRepo := overtimeRepositoryPkg.NewOvertimeRepository()
	payrollRepo := payrollRepositoryPkg.NewPayrollRepository()
	payslipRepo := payslipRepositoryPkg.NewPayslipRepository()
	reimbursementRepo := reimbursementRepositoryPkg.NewReimbursementRepository()

	// Services
	attendanceService := attendanceServicePkg.NewAttendanceService(attendanceRepo)
	overtimeService := overtimeServicePkg.NewOvertimeService(overtimeRepo)
	payrollService := payrollServicePkg.NewPayrollService(payrollRepo)
	payslipService := payslipServicePkg.NewPayslipService(payslipRepo)
	reimbursementService := reimbursementServicePkg.NewReimbursementService(reimbursementRepo)

	// Handler
	attendanceHandler := &attendanceHandlerPkg.AttendanceHandler{
		Service: attendanceService,
	}
	overtimeHandler := &overtimeHandlerPkg.OvertimeHandler{
		Service: overtimeService,
	}
	payrollHandler := &payrollHandlerPkg.PayrollHandler{
		Service:       payrollService,
		GetPeriodByID: payrollRepositoryPkg.NewPayrollRepository().GetPeriodByID,
	}
	payslipHandler := &payslipHandlerPkg.PayslipHandler{
		Service: payslipService,
	}
	reimbursementHandler := &reimbursementHandlerPkg.ReimbursementHandler{
		Service: reimbursementService,
	}

	// Public route
	e.POST("/login", handlers.Login)

	auth := e.Group("/api", middleware.JWTAuth) // Re-apply JWT auth middleware
	// Employee routes:
	auth.GET("/me", handlers.GetMyProfile)
	auth.GET("/attendance", attendanceHandler.GetMyAttendance, middleware.RequireRole("employee"))
	auth.POST("/attendance", attendanceHandler.SubmitAttendance, middleware.RequireRole("employee"))
	auth.GET("/overtime", overtimeHandler.GetMyOvertime, middleware.RequireRole("employee"))
	auth.POST("/overtime", overtimeHandler.SubmitOvertime, middleware.RequireRole("employee"))
	auth.GET("/reimbursement", reimbursementHandler.GetMyReimbursements, middleware.RequireRole("employee"))
	auth.POST("/reimbursement", reimbursementHandler.SubmitReimbursement, middleware.RequireRole("employee"))

	auth.GET("/payslip/:period_id", payslipHandler.GetMyPayslip, middleware.RequireRole("employee"))

	// Admin routes
	auth.POST("/payroll-period", payrollHandler.CreatePayrollPeriod, middleware.RequireRole("admin"))
	auth.POST("/run-payroll", payrollHandler.RunPayroll, middleware.RequireRole("admin"))
	auth.GET("/payslip-summary/:period_id", payslipHandler.GetPayslipSummary, middleware.RequireRole("admin"))

}
