package database

import (
	"api-page/main/src/models"

	"gorm.io/gorm"
)

// Migrate the database schema.
// See: https://gorm.io/docs/migration.html#Auto-Migration
func Migrate(db *gorm.DB) error {
	// Updated migration set: normalized models + existing domain models.
	err := db.AutoMigrate(&models.App{})
	if err != nil {
		return err
	}

	return nil
}
