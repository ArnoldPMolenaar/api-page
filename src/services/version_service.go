package services

import (
	"api-page/main/src/cache"
	"api-page/main/src/database"
	"api-page/main/src/dto/requests"
	"api-page/main/src/dto/responses"
	"api-page/main/src/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/gofiber/fiber/v2"
	"github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
)

// IsVersionPublished method to check if a version is published.
func IsVersionPublished(versionID uint) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.Version{}, "id = ? AND published_at IS NOT NULL", versionID); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// IsVersionAvailable method to check if a version is available.
func IsVersionAvailable(appName, versionName string, ignore *string) (bool, error) {
	query := database.Pg.Limit(1)
	var result *gorm.DB
	if ignore != nil {
		result = query.Find(&models.Version{}, "app_name = ? AND name = ? AND name != ?", appName, versionName, ignore)
	} else {
		result = query.Find(&models.Version{}, "app_name = ? AND name = ?", appName, versionName)
	}

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected == 0, nil
}

// IsVersionDeleted method to check if a version is deleted.
func IsVersionDeleted(versionID uint) (bool, error) {
	if result := database.Pg.Unscoped().Limit(1).Find(&models.Version{}, "id = ? AND deleted_at IS NOT NULL", versionID); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// GetVersions method to get paginated versions.
func GetVersions(c *fiber.Ctx) (*pagination.Model, error) {
	versions := make([]models.Version, 0)
	values := c.Request().URI().QueryArgs()
	allowedColumns := map[string]bool{
		"id":           true,
		"publish_id":   true,
		"app_name":     true,
		"name":         true,
		"enabled_at":   true,
		"published_at": true,
		"created_at":   true,
		"updated_at":   true,
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
	dbResult := database.Pg.Scopes(queryFunc, sortFunc).
		Limit(limit).
		Offset(offset)

	total := int64(0)
	dbCount := database.Pg.Scopes(queryFunc).
		Model(&models.Version{})

	if result := dbResult.Find(&versions); result.Error != nil {
		return nil, result.Error
	}

	dbCount.Count(&total)
	pageCount := pagination.Count(int(total), limit)

	paginatedVersions := make([]responses.PaginatedVersion, 0)
	for i := range versions {
		paginatedVersion := responses.PaginatedVersion{}
		paginatedVersion.SetPaginatedVersion(&versions[i])
		paginatedVersions = append(paginatedVersions, paginatedVersion)
	}

	paginationModel := pagination.CreatePaginationModel(limit, page, pageCount, int(total), paginatedVersions)

	return &paginationModel, nil
}

// GetVersionLookup method to get a lookup of versions.
func GetVersionLookup(appName string, name *string) (*[]models.Version, error) {
	versions := make([]models.Version, 0)

	if inCache, err := isVersionsLookupInCache(appName); err != nil {
		return nil, err
	} else if inCache {
		if cacheVersions, err := getVersionsLookupFromCache(appName); err != nil {
			return nil, err
		} else if cacheVersions != nil && len(*cacheVersions) > 0 {
			versions = *cacheVersions
		}
	}

	if len(versions) == 0 {
		query := database.Pg.Model(&models.Version{}).
			Select("id", "name")

		if result := query.Find(&versions, "app_name = ? AND enabled_at IS NOT NULL", appName); result.Error != nil {
			return nil, result.Error
		}

		_ = setVersionsLookupToCache(appName, &versions)
	}

	// If a name filter is provided, perform case-insensitive substring match on the list.
	if name != nil {
		target := strings.TrimSpace(*name)
		if target != "" {
			lowerTarget := strings.ToLower(target)
			filtered := make([]models.Version, 0, len(versions))
			for i := range versions {
				if strings.Contains(strings.ToLower(versions[i].Name), lowerTarget) {
					filtered = append(filtered, versions[i])
				}
			}
			versions = filtered
		}
	}

	return &versions, nil
}

// GetVersionByID method to get a version by ID.
func GetVersionByID(versionID uint) (*models.Version, error) {
	version := &models.Version{}

	if result := database.Pg.Find(version, "id = ?", versionID); result.Error != nil {
		return nil, result.Error
	}

	return version, nil
}

// GetPublishedVersionByAppName method to get the published version by app name.
func GetPublishedVersionByAppName(appName string) (*models.Version, error) {
	version := &models.Version{}

	if result := database.Pg.Limit(1).Find(version, "app_name = ? AND published_at IS NOT NULL", appName); result.Error != nil {
		return nil, result.Error
	}

	return version, nil
}

// CreateVersion method to create a version.
func CreateVersion(version *requests.CreateVersion) (*models.Version, error) {
	v := &models.Version{AppName: version.AppName, Name: version.Name}
	if version.EnabledAt != nil {
		v.EnabledAt = sql.NullTime{Time: *version.EnabledAt, Valid: true}
	}

	result := &models.Version{}
	if err := database.Pg.FirstOrCreate(&result, v).Error; err != nil {
		return nil, err
	}

	_ = deleteVersionsLookupFromCache(v.AppName)

	return result, nil
}

// UpdateVersion method to update a version.
func UpdateVersion(oldVersion models.Version, version *requests.UpdateVersion) (*models.Version, error) {
	oldVersion.Name = version.Name
	if version.EnabledAt != nil {
		oldVersion.EnabledAt = sql.NullTime{Time: *version.EnabledAt, Valid: true}
	} else {
		oldVersion.EnabledAt = sql.NullTime{Valid: false}
	}

	if result := database.Pg.Save(&oldVersion); result.Error != nil {
		return nil, result.Error
	}

	_ = deleteVersionsLookupFromCache(oldVersion.AppName)

	return &oldVersion, nil
}

// PublishVersion method to publish a version.
func PublishVersion(appName string, versionID uint) error {
	// Start a new transaction
	tx := database.Pg.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if result := tx.Model(&models.Version{}).Where("app_name = ?", appName).Update("published_at", sql.NullTime{Valid: false}); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result := tx.Model(&models.Version{}).Where("id = ?", versionID).Update("published_at", sql.NullTime{Time: time.Now(), Valid: true}); result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// DeleteVersion method to delete a version.
func DeleteVersion(versionID uint, appName string) error {
	err := database.Pg.Delete(&models.Version{}, versionID).Error
	if err == nil {
		_ = deleteVersionsLookupFromCache(appName)
	}

	return err
}

// RestoreVersion method to restore a deleted version.
func RestoreVersion(versionID uint) error {
	err := database.Pg.Unscoped().Model(&models.Version{}).Where("id = ?", versionID).Update("deleted_at", nil).Error
	if err == nil {
		var appName string

		if result := database.Pg.Unscoped().Model(&models.Version{}).Where("id = ?", versionID).Pluck("app_name", &appName); result.Error != nil {
			return result.Error
		}

		_ = deleteVersionsLookupFromCache(appName)
	}

	return err
}

// getVersionsLookupCacheKey gets the key for the cache.
func getVersionsLookupCacheKey(appName string) string {
	return fmt.Sprintf("versions:lookup:%s", appName)
}

// isVersionsLookupInCache checks if the versions exists in the cache.
func isVersionsLookupInCache(appName string) (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(getVersionsLookupCacheKey(appName)).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// getVersionsLookupFromCache gets the versions from the cache.
func getVersionsLookupFromCache(appName string) (*[]models.Version, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(getVersionsLookupCacheKey(appName)).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var versions []models.Version
	if err := json.Unmarshal([]byte(value), &versions); err != nil {
		return nil, err
	}

	return &versions, nil
}

// setVersionsLookupToCache sets the versions to the cache.
func setVersionsLookupToCache(appName string, versions *[]models.Version) error {
	value, err := json.Marshal(versions)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(getVersionsLookupCacheKey(appName)).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deleteVersionsLookupFromCache deletes existing versions from the cache.
func deleteVersionsLookupFromCache(appName string) error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(getVersionsLookupCacheKey(appName)).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}
