package repositories

import (
	"school_attendance_backend/internal/models"

	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

func (r *StudentRepository) Update(id string, student *models.Student) error {
	return r.db.Model(&models.Student{}).Where("id = ?", id).Updates(student).Error
}

func (r *StudentRepository) Delete(id string) error {
	return r.db.Delete(&models.Student{}, "id = ?", id).Error
}

func (r *StudentRepository) GetByID(id string) (*models.Student, error) {
	var student models.Student
	err := r.db.First(&student, "id = ?", id).Error
	return &student, err
}

func (r *StudentRepository) GetByStudentID(studentID string) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("student_id = ?", studentID).First(&student).Error
	return &student, err
}

func (r *StudentRepository) GetAll() ([]models.Student, error) {
	var students []models.Student
	err := r.db.Find(&students).Error
	return students, err
}
