package services

import (
	"api-page/main/src/database"
	"api-page/main/src/dto/responses"
	"api-page/main/src/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// IsAppAvailable checks whether an app with the given name exists.
// It returns true when a matching row is found.
func IsAppAvailable(app string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.App{}, "name = ?", app); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// GetApps returns all apps from the database.
func GetApps() (*[]models.App, error) {
	apps := make([]models.App, 0)

	if result := database.Pg.Find(&apps); result.Error != nil {
		return nil, result.Error
	}

	return &apps, nil
}

// GetAppModuleTypes fetches the module types associated with an app.
// The preload join excludes soft-deleted module types.
func GetAppModuleTypes(app string) ([]models.ModuleType, error) {
	a := &models.App{}

	if result := database.Pg.Preload("ModuleTypes").Find(a, "name = ?", app); result.Error != nil {
		return nil, result.Error
	}

	return a.ModuleTypes, nil
}

// GetAppPluginTypes fetches the plugin types associated with an app.
func GetAppPluginTypes(app string) ([]models.PluginType, error) {
	a := &models.App{}

	if result := database.Pg.Preload("PluginTypes").Find(a, "name = ?", app); result.Error != nil {
		return nil, result.Error
	}

	return a.PluginTypes, nil
}

// CreateApp creates the app when it does not exist yet.
// If the app already exists, the existing row is reused.
func CreateApp(name string) (*models.App, error) {
	app := &models.App{Name: name}

	if err := database.Pg.FirstOrCreate(&models.App{}, app).Error; err != nil {
		return nil, err
	}

	return app, nil
}

// dedupeStrings removes duplicate values while preserving first-seen order.
func dedupeStrings(values []string) []string {
	unique := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		unique = append(unique, value)
	}

	return unique
}

// syncAppAssociationByName synchronizes a named app association to match
// exactly the provided desired names: missing links are created and stale
// links are removed, while unchanged links are kept.
//
// The function is generic so it can be reused for module types, plugin types,
// and any future app association keyed by a name field.
func syncAppAssociationByName[T any](
	tx *gorm.DB,
	appName string,
	associationName string,
	desiredNames []string,
	entityLabel string,
	loadDesired func([]string) ([]T, error),
	getName func(T) string,
) error {
	a := models.App{}
	if err := tx.First(&a, "name = ?", appName).Error; err != nil {
		return err
	}

	dedupedDesiredNames := dedupeStrings(desiredNames)
	desiredEntities := make([]T, 0, len(dedupedDesiredNames))
	if len(dedupedDesiredNames) > 0 {
		loaded, err := loadDesired(dedupedDesiredNames)
		if err != nil {
			return err
		}
		desiredEntities = loaded

		if len(desiredEntities) != len(dedupedDesiredNames) {
			found := make(map[string]struct{}, len(desiredEntities))
			for _, entity := range desiredEntities {
				found[getName(entity)] = struct{}{}
			}

			missing := make([]string, 0)
			for _, name := range dedupedDesiredNames {
				if _, ok := found[name]; !ok {
					missing = append(missing, name)
				}
			}

			return fmt.Errorf("%s(s) not found: %s", entityLabel, strings.Join(missing, ", "))
		}
	}

	currentEntities := make([]T, 0)
	association := tx.Model(&a).Association(associationName)
	if err := association.Find(&currentEntities); err != nil {
		return err
	}

	desiredNamesSet := make(map[string]struct{}, len(desiredEntities))
	for _, entity := range desiredEntities {
		desiredNamesSet[getName(entity)] = struct{}{}
	}

	currentNamesSet := make(map[string]struct{}, len(currentEntities))
	for _, entity := range currentEntities {
		currentNamesSet[getName(entity)] = struct{}{}
	}

	toDelete := make([]T, 0)
	for _, entity := range currentEntities {
		if _, ok := desiredNamesSet[getName(entity)]; !ok {
			toDelete = append(toDelete, entity)
		}
	}

	toCreate := make([]T, 0)
	for _, entity := range desiredEntities {
		if _, ok := currentNamesSet[getName(entity)]; !ok {
			toCreate = append(toCreate, entity)
		}
	}

	if len(toDelete) > 0 {
		if err := association.Delete(&toDelete); err != nil {
			return err
		}
	}

	if len(toCreate) > 0 {
		if err := association.Append(&toCreate); err != nil {
			return err
		}
	}

	return nil
}

// SetAppModuleTypes synchronizes app -> module_type associations.
// It validates all requested module types exist and then applies a delta sync.
func SetAppModuleTypes(app string, moduleTypes []string) (*responses.AppTypes, error) {
	var response *responses.AppTypes

	err := database.Pg.Transaction(func(tx *gorm.DB) error {
		if err := syncAppAssociationByName(
			tx,
			app,
			"ModuleTypes",
			moduleTypes,
			"module type",
			func(names []string) ([]models.ModuleType, error) {
				desiredModuleTypes := make([]models.ModuleType, 0, len(names))
				if err := tx.Where("name IN ?", names).Find(&desiredModuleTypes).Error; err != nil {
					return nil, err
				}

				return desiredModuleTypes, nil
			},
			func(entity models.ModuleType) string {
				return entity.Name
			},
		); err != nil {
			return err
		}

		types := make([]models.ModuleType, 0)
		if err := tx.Model(&models.App{Name: app}).Association("ModuleTypes").Find(&types); err != nil {
			return err
		}

		response = &responses.AppTypes{}
		response.SetAppModuleTypes(app, types)

		return nil
	})
	if err != nil {
		return nil, err
	}

	_ = deleteModuleTypesLookupFromCache(nil)
	_ = deleteModuleTypesLookupFromCache(&app)

	return response, nil
}

// SetAppPluginTypes synchronizes app -> plugin_type associations.
// It validates all requested plugin types exist and then applies a delta sync.
func SetAppPluginTypes(app string, pluginTypes []string) (*responses.AppTypes, error) {
	var response *responses.AppTypes

	err := database.Pg.Transaction(func(tx *gorm.DB) error {
		if err := syncAppAssociationByName(
			tx,
			app,
			"PluginTypes",
			pluginTypes,
			"plugin type",
			func(names []string) ([]models.PluginType, error) {
				desiredPluginTypes := make([]models.PluginType, 0, len(names))
				if err := tx.Where("name IN ?", names).Find(&desiredPluginTypes).Error; err != nil {
					return nil, err
				}

				return desiredPluginTypes, nil
			},
			func(entity models.PluginType) string {
				return entity.Name
			},
		); err != nil {
			return err
		}

		types := make([]models.PluginType, 0)
		if err := tx.Model(&models.App{Name: app}).Association("PluginTypes").Find(&types); err != nil {
			return err
		}

		response = &responses.AppTypes{}
		response.SetAppPluginTypes(app, types)

		return nil
	})
	if err != nil {
		return nil, err
	}

	_ = deletePluginTypesLookupFromCache(nil)
	_ = deletePluginTypesLookupFromCache(&app)

	return response, nil
}
