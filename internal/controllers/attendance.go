package controllers

import (
	"school_attendance_backend/internal/models"
	"school_attendance_backend/internal/repositories"
	"school_attendance_backend/internal/services"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type AttendanceController struct {
	attendanceRepo *repositories.AttendanceRepository
	studentRepo    *repositories.StudentRepository
	emailService   *services.EmailService
	qrcodeService  *services.QRCodeService
}

func NewAttendanceController(
	attendanceRepo *repositories.AttendanceRepository,
	studentRepo *repositories.StudentRepository,
	emailService *services.EmailService,
	qrcodeService *services.QRCodeService,
) *AttendanceController {
	return &AttendanceController{
		attendanceRepo: attendanceRepo,
		studentRepo:    studentRepo,
		emailService:   emailService,
		qrcodeService:  qrcodeService,
	}
}

func (c *AttendanceController) GetAllAttendance(ctx fiber.Ctx) error {
	attendances, err := c.attendanceRepo.GetAllAttendance()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendance records",
		})
	}
	return ctx.JSON(attendances)
}

func (c *AttendanceController) RecordAttendance(ctx fiber.Ctx) error {
	var req models.AttendanceRecordRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Parse QR code data
	qrData, err := c.qrcodeService.Parse(req.QRCodeData)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid QR code data",
		})
	}

	studentID, ok := qrData["student_id"]
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing student ID in QR code",
		})
	}

	// Get student
	student, err := c.studentRepo.GetByStudentID(studentID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	now := time.Now()
	cutoffTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, now.Location())
	status := "present"
	if now.After(cutoffTime) {
		status = "late"
	}

	attendance := models.Attendance{
		ID:               uuid.New().String(),
		StudentID:        student.ID,
		StudentName:      student.Name,
		ParentEmail:      student.ParentEmail,
		Timestamp:        now,
		Status:           status,
		NotificationSent: false,
	}

	// Send arrival notification
	if err := c.emailService.SendArrivalNotification(
		student.ParentEmail,
		student.Name,
		now.Format("15:04"),
	); err == nil {
		attendance.NotificationSent = true
	}

	if err := c.attendanceRepo.Create(&attendance); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to record attendance",
		})
	}

	// Update student's last attendance
	student.LastAttendance = &now
	if err := c.studentRepo.Update(student.ID, student); err != nil {
		// Log error but don't fail the request
	}

	return ctx.Status(fiber.StatusCreated).JSON(attendance)
}

func (c *AttendanceController) GetAttendanceByDate(ctx fiber.Ctx) error {
	dateStr := ctx.Params("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format (YYYY-MM-DD)",
		})
	}

	attendances, err := c.attendanceRepo.GetByDate(date)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get attendance records",
		})
	}

	return ctx.JSON(attendances)
}

func (c *AttendanceController) GetAbsentStudents(ctx fiber.Ctx) error {
	dateStr := ctx.Params("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format (YYYY-MM-DD)",
		})
	}

	absentStudents, err := c.attendanceRepo.GetAbsentStudents(date)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get absent students",
		})
	}

	return ctx.JSON(absentStudents)
}

func (c *AttendanceController) SendAbsenceNotifications(ctx fiber.Ctx) error {
	dateStr := ctx.Params("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format (YYYY-MM-DD)",
		})
	}

	absentStudents, err := c.attendanceRepo.GetAbsentStudents(date)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get absent students",
		})
	}

	// Get student details for absent students
	var students []models.Student
	for _, studentID := range absentStudents {
		student, err := c.studentRepo.GetByStudentID(studentID)
		if err != nil {
			continue // Skip if student not found
		}
		students = append(students, *student)
	}

	// Send notifications
	for _, student := range students {
		if err := c.emailService.SendAbsenceNotification(
			student.ParentEmail,
			student.Name,
			date.Format("2006-01-02"),
		); err == nil {
			// Record notification in attendance
			attendance := models.Attendance{
				StudentID:        student.ID,
				StudentName:      student.Name,
				ParentEmail:      student.ParentEmail,
				Timestamp:        time.Now(),
				Status:           "absent",
				NotificationSent: true,
			}
			_ = c.attendanceRepo.Create(&attendance)
		}
	}

	return ctx.JSON(fiber.Map{
		"message": "Absence notifications processed",
		"count":   len(students),
	})
}
