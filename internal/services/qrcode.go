package services

import (
	"encoding/json"
	"school_attendance_backend/internal/models"
	"time"
)

type QRCodeService struct{}

func NewQRCodeService() *QRCodeService {
	return &QRCodeService{}
}

func (s *QRCodeService) Generate(student *models.Student) (string, error) {
	data := map[string]interface{}{
		"student_id":   student.StudentID,
		"student_name": student.Name,
		"parent_email": student.ParentEmail,
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (s *QRCodeService) Parse(qrCodeData string) (map[string]string, error) {
	var data map[string]string
	err := json.Unmarshal([]byte(qrCodeData), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
