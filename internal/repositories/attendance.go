package repositories

import (
	"school_attendance_backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) Create(attendance *models.Attendance) error {
	return r.db.Create(attendance).Error
}

func (r *AttendanceRepository) GetAllAttendance() ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := r.db.Find(&attendances).Error
	return attendances, err
}

func (r *AttendanceRepository) GetByStudentID(studentID string) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := r.db.Where("student_id = ?", studentID).Find(&attendances).Error
	return attendances, err
}

func (r *AttendanceRepository) GetByDate(date time.Time) ([]models.Attendance, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1)

	var attendances []models.Attendance
	err := r.db.Where("timestamp >= ? AND timestamp < ?", start, end).Find(&attendances).Error
	return attendances, err
}

func (r *AttendanceRepository) GetAbsentStudents(date time.Time) ([]string, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1)

	var presentStudentIDs []string
	err := r.db.Model(&models.Attendance{}).
		Where("timestamp >= ? AND timestamp < ?", start, end).
		Distinct("student_id").
		Pluck("student_id", &presentStudentIDs).Error
	if err != nil {
		return nil, err
	}

	var absentStudentIDs []string
	err = r.db.Model(&models.Student{}).
		Where("id NOT IN ?", presentStudentIDs).
		Pluck("student_id", &absentStudentIDs).Error

	return absentStudentIDs, err
}
