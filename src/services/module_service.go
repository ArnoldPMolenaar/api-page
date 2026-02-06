package services

import (
	"api-page/main/src/cache"
	"api-page/main/src/database"
	"api-page/main/src/dto/requests"
	"api-page/main/src/dto/responses"
	"api-page/main/src/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/gofiber/fiber/v2"
	"github.com/valkey-io/valkey-go"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// IsModuleNameAvailable method to check if a module is available.
func IsModuleNameAvailable(appName, name string, ignore *string) (bool, error) {
	query := database.Pg.Limit(1)
	var result *gorm.DB
	if ignore != nil {
		result = query.Find(&models.Module{}, "app_name = ? AND name = ? AND name != ?", appName, name, ignore)
	} else {
		result = query.Find(&models.Module{}, "app_name = ? AND name = ?", appName, name)
	}

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected == 0, nil
}

// IsModuleTypeNotAvailable method to check if a module type is available.
func IsModuleTypeNotAvailable(moduleType string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.ModuleType{}, "name = ?", moduleType); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 0, nil
	}
}

// IsModuleDeleted method to check if a module is deleted.
func IsModuleDeleted(moduleID uint) (bool, error) {
	if result := database.Pg.Unscoped().Limit(1).Find(&models.Module{}, "id = ? AND deleted_at IS NOT NULL", moduleID); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// GetModules method to get paginated modules.
func GetModules(c *fiber.Ctx) (*pagination.Model, error) {
	modules := make([]models.Module, 0)
	values := c.Request().URI().QueryArgs()
	allowedColumns := map[string]bool{
		"id":         true,
		"app_name":   true,
		"type":       true,
		"name":       true,
		"created_at": true,
		"updated_at": true,
	}

	queryFunc := pagination.Query(values, allowedColumns)
	sortFunc := pagination.Sort(values, allowedColumns)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}
	offset := pagination.Offset(page, limit)
	dbResult := database.Pg.Scopes(queryFunc, sortFunc, scopeExcludeDeletedModuleType).
		Limit(limit).
		Offset(offset)

	total := int64(0)
	dbCount := database.Pg.Scopes(queryFunc, scopeExcludeDeletedModuleType).
		Model(&models.Module{})

	if result := dbResult.Find(&modules); result.Error != nil {
		return nil, result.Error
	}

	dbCount.Count(&total)
	pageCount := pagination.Count(int(total), limit)

	paginatedModules := make([]responses.PaginatedModule, 0)
	for i := range modules {
		paginatedModule := responses.PaginatedModule{}
		paginatedModule.SetPaginatedModule(&modules[i])
		paginatedModules = append(paginatedModules, paginatedModule)
	}

	paginationModel := pagination.CreatePaginationModel(limit, page, pageCount, int(total), paginatedModules)

	return &paginationModel, nil
}

// GetModuleLookup method to get a lookup of modules.
func GetModuleLookup(appName string, name *string) (*[]models.Module, error) {
	modules := make([]models.Module, 0)

	if inCache, err := isModulesLookupInCache(appName); err != nil {
		return nil, err
	} else if inCache {
		if cacheModules, err := getModulesLookupFromCache(appName); err != nil {
			return nil, err
		} else if cacheModules != nil && len(*cacheModules) > 0 {
			modules = *cacheModules
		}
	}

	if len(modules) == 0 {
		query := database.Pg.Model(&models.Module{}).
			Scopes(scopeExcludeDeletedModuleType).
			Select("modules.id", "modules.name")

		if result := query.Find(&modules, "app_name = ?", appName); result.Error != nil {
			return nil, result.Error
		}

		_ = setModulesLookupToCache(appName, &modules)
	}

	// If a name filter is provided, perform case-insensitive substring match on the list.
	if name != nil {
		target := strings.TrimSpace(*name)
		if target != "" {
			lowerTarget := strings.ToLower(target)
			filtered := make([]models.Module, 0, len(modules))
			for i := range modules {
				if strings.Contains(strings.ToLower(modules[i].Name), lowerTarget) {
					filtered = append(filtered, modules[i])
				}
			}
			modules = filtered
		}
	}

	return &modules, nil
}

// GetModuleByID method to get a module by ID.
func GetModuleByID(moduleID uint) (*models.Module, error) {
	module := &models.Module{}

	if result := database.Pg.Scopes(scopeExcludeDeletedModuleType).Find(module, "id = ?", moduleID); result.Error != nil {
		return nil, result.Error
	}

	return module, nil
}

// CreateModule method to create a module.
func CreateModule(module *requests.CreateModule) (*models.Module, error) {
	m := &models.Module{
		AppName:  module.AppName,
		Type:     module.Type,
		Name:     module.Name,
		Settings: datatypes.JSON(module.Settings),
	}

	result := &models.Module{}
	if err := database.Pg.FirstOrCreate(&result, m).Error; err != nil {
		return nil, err
	}

	_ = deleteModulesLookupFromCache(m.AppName)

	return result, nil
}

// UpdateModule method to update a module.
func UpdateModule(oldModule models.Module, module *requests.UpdateModule) (*models.Module, error) {
	oldModule.Name = module.Name
	oldModule.Type = module.Type
	oldModule.Settings = datatypes.JSON(module.Settings)

	if result := database.Pg.Save(&oldModule); result.Error != nil {
		return nil, result.Error
	}

	_ = deleteModulesLookupFromCache(oldModule.AppName)

	return &oldModule, nil
}

// DeleteModule method to delete a module.
func DeleteModule(moduleID uint, appName string) error {
	err := database.Pg.Delete(&models.Module{}, moduleID).Error
	if err == nil {
		_ = deleteModulesLookupFromCache(appName)
	}

	return err
}

// RestoreModule method to restore a deleted module.
func RestoreModule(moduleID uint) error {
	err := database.Pg.Unscoped().Model(&models.Module{}).Where("id = ?", moduleID).Update("deleted_at", nil).Error
	if err == nil {
		var appName string

		if result := database.Pg.Unscoped().Model(&models.Module{}).Where("id = ?", moduleID).Pluck("app_name", &appName); result.Error != nil {
			return result.Error
		}

		_ = deleteModulesLookupFromCache(appName)
	}

	return err
}

// getModulesLookupCacheKey gets the key for the cache.
func getModulesLookupCacheKey(appName string) string {
	return fmt.Sprintf("modules:lookup:%s", appName)
}

// isModulesLookupInCache checks if the modules exists in the cache.
func isModulesLookupInCache(appName string) (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(getModulesLookupCacheKey(appName)).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// getModulesLookupFromCache gets the modules from the cache.
func getModulesLookupFromCache(appName string) (*[]models.Module, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(getModulesLookupCacheKey(appName)).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var modules []models.Module
	if err := json.Unmarshal([]byte(value), &modules); err != nil {
		return nil, err
	}

	return &modules, nil
}

// setModulesLookupToCache sets the modules to the cache.
func setModulesLookupToCache(appName string, modules *[]models.Module) error {
	value, err := json.Marshal(modules)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(getModulesLookupCacheKey(appName)).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deleteModulesLookupFromCache deletes existing modules from the cache.
func deleteModulesLookupFromCache(appName string) error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(getModulesLookupCacheKey(appName)).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// scopeExcludeDeletedModuleType excludes modules whose Type was soft-deleted.
func scopeExcludeDeletedModuleType(db *gorm.DB) *gorm.DB {
	return db.
		Joins("JOIN module_types ON module_types.name = modules.type").
		Where("module_types.deleted_at IS NULL")
}
