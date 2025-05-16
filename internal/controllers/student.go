package controllers

import (
	"school_attendance_backend/internal/models"
	"school_attendance_backend/internal/repositories"
	"school_attendance_backend/internal/services"
	"school_attendance_backend/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type StudentController struct {
	studentRepo   *repositories.StudentRepository
	qrcodeService *services.QRCodeService
}

func NewStudentController(
	studentRepo *repositories.StudentRepository,
	qrcodeService *services.QRCodeService,
) *StudentController {
	return &StudentController{
		studentRepo:   studentRepo,
		qrcodeService: qrcodeService,
	}
}

func (c *StudentController) CreateStudent(ctx fiber.Ctx) error {
	var req models.CreateStudentRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate QR code
	student := models.Student{
		ID:          req.ID,
		Name:        req.Name,
		StudentID:   req.StudentID,
		Email:       req.Email,
		Grade:       req.Grade,
		ParentEmail: req.ParentEmail,
		ParentPhone: req.ParentPhone,
	}

	qrCode, err := c.qrcodeService.Generate(&student)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate QR code",
		})
	}
	student.QRCodeValue = qrCode

	// Create student
	if err := c.studentRepo.Create(&student); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create student",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(student)
}

func (c *StudentController) GetByID(ctx fiber.Ctx) error {
	studentID := ctx.Params("id")

	student, err := c.studentRepo.GetByID(studentID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	return ctx.JSON(student)
}

func (c *StudentController) UpdateStudent(ctx fiber.Ctx) error {
	studentID := ctx.Params("id")
	if studentID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is required",
		})
	}

	var req models.UpdateStudentRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Update student
	student := &models.Student{
		Name:        req.Name,
		StudentID:   studentID,
		Email:       req.Email,
		Grade:       req.Grade,
		ParentEmail: req.ParentEmail,
		ParentPhone: req.ParentPhone,
	}
	if err := c.studentRepo.Update(studentID, student); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update student",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Student updated successfully",
	})
}

func (c *StudentController) DeleteStudent(ctx fiber.Ctx) error {
	studentID := ctx.Params("id")
	if studentID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is required",
		})
	}

	if err := c.studentRepo.Delete(studentID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete student",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "Student deleted successfully",
	})
}

func (c *StudentController) GetAllStudents(ctx fiber.Ctx) error {
	students, err := c.studentRepo.GetAll()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve students",
		})
	}

	return ctx.JSON(students)
}

func (c *StudentController) GetByStudentID(ctx fiber.Ctx) error {
	studentID := ctx.Params("student_id")
	if studentID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is required",
		})
	}

	student, err := c.studentRepo.GetByStudentID(studentID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Student not found",
		})
	}

	return ctx.JSON(student)
}
