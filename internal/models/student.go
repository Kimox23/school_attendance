package models

import "time"

type Student struct {
	ID             string     `json:"id" gorm:"primaryKey"`
	Name           string     `json:"name" gorm:"not null"`
	StudentID      string     `json:"student_id" gorm:"unique;not null"`
	Email          string     `json:"email" gorm:"not null"`
	Grade          string     `json:"grade" gorm:"not null"`
	ParentEmail    string     `json:"parent_email" gorm:"not null"`
	ParentPhone    string     `json:"parent_phone"`
	QRCodeValue    string     `json:"qr_code_value" gorm:"not null"`
	PhotoURL       string     `json:"photo_url"`
	LastAttendance *time.Time `json:"last_attendance"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateStudentRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name" validate:"required"`
	StudentID   string `json:"student_id" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Grade       string `json:"grade" validate:"required"`
	ParentEmail string `json:"parent_email" validate:"required,email"`
	ParentPhone string `json:"parent_phone"`
}

type UpdateStudentRequest struct {
	Name        string `json:"name"`
	StudentID   string `json:"student_id"`
	Email       string `json:"email" validate:"omitempty,email"`
	Grade       string `json:"grade"`
	ParentEmail string `json:"parent_email" validate:"omitempty,email"`
	ParentPhone string `json:"parent_phone"`
}
