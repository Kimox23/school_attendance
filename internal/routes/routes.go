package routes

import (
	"school_attendance_backend/internal/config"
	"school_attendance_backend/internal/controllers"
	"school_attendance_backend/internal/middleware"
	"school_attendance_backend/internal/repositories"
	"school_attendance_backend/internal/services"
	"school_attendance_backend/internal/utils"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, config *config.Config) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	studentRepo := repositories.NewStudentRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)

	// Initialize services
	emailService := services.NewEmailService(config.ResendAPIKey, config.SenderEmail)
	qrcodeService := services.NewQRCodeService()
	jwtUtil := utils.NewJWTUtil(config.JWTSecret, config.JWTExpiration)

	// Initialize controllers
	authController := controllers.NewAuthController(userRepo, jwtUtil)
	studentController := controllers.NewStudentController(studentRepo, qrcodeService)
	attendanceController := controllers.NewAttendanceController(
		attendanceRepo,
		studentRepo,
		emailService,
		qrcodeService,
	)

	// Auth routes
	auth := app.Group("/api/auth")
	{
		auth.Post("/register", authController.Register)
		auth.Post("/login", authController.Login)
	}

	// Student routes (protected)
	student := app.Group("/api/students", middleware.AuthRequired(jwtUtil))
	{
		student.Post("/", studentController.CreateStudent)
		student.Get("/", studentController.GetAllStudents)
		student.Get("/id/:id", studentController.GetByID)
		student.Get("/student_id/:student_id", studentController.GetByStudentID)
		student.Put("/:id", studentController.UpdateStudent)
		student.Delete("/:id", studentController.DeleteStudent)
	}

	// Attendance routes (protected)
	attendance := app.Group("/api/attendance", middleware.AuthRequired(jwtUtil))
	{
		attendance.Get("/", attendanceController.GetAllAttendance)
		attendance.Post("/", attendanceController.RecordAttendance)
		attendance.Get("/date/:date", attendanceController.GetAttendanceByDate)
		attendance.Get("/absent/:date", attendanceController.GetAbsentStudents)
		attendance.Post("/notify-absent/:date", attendanceController.SendAbsenceNotifications)
	}
}
