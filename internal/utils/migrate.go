package utils

import (
	"school_attendance_backend/internal/models"

	"gorm.io/gorm"
)

func MigrateModels(db *gorm.DB) error {
	// Disable foreign key checks temporarily
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		return err
	}

	// List all your models here
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Student{},
		&models.Attendance{},
	}

	// Perform auto-migration for each model
	for _, model := range modelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			// Re-enable foreign key checks if migration fails
			_ = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
			return err
		}
	}

	// Re-enable foreign key checks
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		return err
	}

	return nil
}
