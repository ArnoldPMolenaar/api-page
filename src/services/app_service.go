package services

import (
	"api-page/main/src/database"
	"api-page/main/src/models"
)

// IsAppAvailable method to check if an app is available.
func IsAppAvailable(app string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.App{}, "name = ?", app); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// GetApps method to get all apps.
func GetApps() (*[]models.App, error) {
	apps := make([]models.App, 0)

	if result := database.Pg.Find(&apps); result.Error != nil {
		return nil, result.Error
	}

	return &apps, nil
}

// CreateApp method to create an app.
func CreateApp(name string) (*models.App, error) {
	app := &models.App{Name: name}

	if err := database.Pg.FirstOrCreate(&models.App{}, app).Error; err != nil {
		return nil, err
	}

	return app, nil
}
