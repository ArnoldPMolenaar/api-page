package services

import (
	"api-page/main/src/cache"
	"api-page/main/src/database"
	"api-page/main/src/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/valkey-io/valkey-go"
)

// GetPluginTypeLookup method to get a lookup of plugin types.
func GetPluginTypeLookup(appName *string) (*[]models.PluginType, error) {
	pluginTypes := make([]models.PluginType, 0)

	if inCache, err := isPluginTypesLookupInCache(appName); err != nil {
		return nil, err
	} else if inCache {
		if cachePluginTypes, err := getPluginTypesLookupFromCache(appName); err != nil {
			return nil, err
		} else if cachePluginTypes != nil && len(*cachePluginTypes) > 0 {
			pluginTypes = *cachePluginTypes
		}
	}

	if len(pluginTypes) == 0 {
		if appName != nil && strings.TrimSpace(*appName) != "" {
			appPluginTypes, err := GetAppPluginTypes(strings.TrimSpace(*appName))
			if err != nil {
				return nil, err
			}
			pluginTypes = appPluginTypes
		} else {
			if result := database.Pg.Find(&pluginTypes); result.Error != nil {
				return nil, result.Error
			}
		}

		_ = setPluginTypesLookupToCache(appName, &pluginTypes)
	}

	return &pluginTypes, nil
}

// getPluginTypesLookupCacheKey gets the key for the cache.
func getPluginTypesLookupCacheKey(appName *string) string {
	if appName == nil || strings.TrimSpace(*appName) == "" {
		return "plugins:types"
	}

	return fmt.Sprintf("plugins:types:%s", strings.TrimSpace(*appName))
}

// isPluginTypesLookupInCache checks if plugin types exist in the cache.
func isPluginTypesLookupInCache(appName *string) (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(getPluginTypesLookupCacheKey(appName)).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// getPluginTypesLookupFromCache gets plugin types from the cache.
func getPluginTypesLookupFromCache(appName *string) (*[]models.PluginType, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(getPluginTypesLookupCacheKey(appName)).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var pluginTypes []models.PluginType
	if err = json.Unmarshal([]byte(value), &pluginTypes); err != nil {
		return nil, err
	}

	return &pluginTypes, nil
}

// setPluginTypesLookupToCache sets plugin types to the cache.
func setPluginTypesLookupToCache(appName *string, pluginTypes *[]models.PluginType) error {
	value, err := json.Marshal(pluginTypes)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(getPluginTypesLookupCacheKey(appName)).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deletePluginTypesLookupFromCache deletes existing plugin type lookups from cache.
func deletePluginTypesLookupFromCache(appName *string) error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(getPluginTypesLookupCacheKey(appName)).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}
