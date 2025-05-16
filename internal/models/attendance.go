package models

import "time"

type Attendance struct {
	ID               string    `json:"id" gorm:"primaryKey"`
	StudentID        string    `json:"student_id" gorm:"not null;index"`
	StudentName      string    `json:"student_name" gorm:"not null"`
	ParentEmail      string    `json:"parent_email" gorm:"not null"`
	Timestamp        time.Time `json:"timestamp" gorm:"not null;index"`
	Status           string    `json:"status" gorm:"not null;type:enum('present','late','absent')"`
	NotificationSent bool      `json:"notification_sent" gorm:"default:false"`
	CreatedAt        time.Time `json:"created_at"`
}

type AttendanceRecordRequest struct {
	QRCodeData string `json:"qr_code_data" validate:"required"`
}

type AttendanceResponse struct {
	ID          string    `json:"id"`
	StudentID   string    `json:"student_id"`
	StudentName string    `json:"student_name"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
}
