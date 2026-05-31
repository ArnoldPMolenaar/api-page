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
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v3"
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
func GetVersions(c fiber.Ctx) (*pagination.Model, error) {
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
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
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
func UpdateVersion(oldVersion *models.Version, version *requests.UpdateVersion) (*models.Version, error) {
	if oldVersion == nil {
		return nil, gorm.ErrRecordNotFound
	}

	oldVersion.Name = version.Name
	if version.EnabledAt != nil {
		oldVersion.EnabledAt = sql.NullTime{Time: *version.EnabledAt, Valid: true}
	} else {
		oldVersion.EnabledAt = sql.NullTime{Valid: false}
	}

	if result := database.Pg.Save(oldVersion); result.Error != nil {
		return nil, result.Error
	}

	_ = deleteVersionsLookupFromCache(oldVersion.AppName)

	return oldVersion, nil
}

// DuplicateVersion method to duplicate a version with selected menus/items and locales,
// including deep cloning of menu trees and pages/partials/indexing, all within a transaction for consistency.
func DuplicateVersion(oldVersion *models.Version, settings *requests.CreateDuplicateVersion) (*models.Version, error) {
	if oldVersion == nil {
		return nil, errors.New("old version is required")
	}

	if settings == nil {
		return nil, errors.New("duplicate version settings are required")
	}

	locales := settings.Locales
	newVersion := &models.Version{AppName: settings.AppName, Name: settings.Name, EnabledAt: utils.NewNullTime(settings.EnabledAt)}
	menuItemMapping := make(map[uint]uint)

	if err := database.Pg.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newVersion).Error; err != nil {
			return err
		}

		sourceMenus := make([]models.Menu, 0)
		if err := tx.
			Preload("MenuItemRelations", func(db *gorm.DB) *gorm.DB {
				return db.
					Preload("MenuItemChild", func(db2 *gorm.DB) *gorm.DB {
						return db2.Preload("Indexing")
					}).
					Order("menu_item_parent_id NULLS FIRST").
					Order("position ASC")
			}).
			Where("version_id = ?", oldVersion.ID).
			Order("id ASC").
			Find(&sourceMenus).Error; err != nil {
			return err
		}

		menuSelections := buildDuplicateMenuSelections(sourceMenus, settings)

		for i := range sourceMenus {
			sourceMenu := &sourceMenus[i]
			selection, ok := menuSelections[sourceMenu.ID]
			if !ok {
				continue
			}

			items, oldPathMap := buildCreateMenuItemsFromSource(sourceMenu, selection)
			if len(items) == 0 {
				continue
			}

			menuCreate := &requests.CreateMenu{
				VersionID: newVersion.ID,
				Name:      sourceMenu.Name,
				Depth:     utils.PtrFromNull[uint8](sourceMenu.Depth),
				Items:     items,
			}

			newMenu, err := CreateMenuWithTx(tx, menuCreate)
			if err != nil {
				return err
			}

			newPathMap := buildMenuItemPathMap(newMenu.MenuItemRelations)

			for path, oldMenuItemID := range oldPathMap {
				if newMenuItemID, ok := newPathMap[path]; ok {
					menuItemMapping[oldMenuItemID] = newMenuItemID
				}
			}
		}

		if err := duplicatePages(tx, menuItemMapping, locales); err != nil {
			return err
		}

		if settings.Footer {
			if err := duplicateFooter(tx, oldVersion.ID, newVersion.ID, locales); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	_ = deleteVersionsLookupFromCache(newVersion.AppName)
	_ = deleteMenusLookupFromCache(newVersion.ID)
	_ = deleteAllVersionMenusFromCache(newVersion.ID)
	for i := range locales {
		_ = deleteVersionMenusFromCache(newVersion.ID, locales[i])
		if settings.Footer {
			_ = deleteFooterFromCache(newVersion.ID, locales[i])
		}
	}
	for _, newMenuItemID := range menuItemMapping {
		for i := range locales {
			_ = deletePageFromCache(newMenuItemID, locales[i])
		}
	}

	return newVersion, nil
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

// duplicateMenuSelection represents the selection of a menu and its items to duplicate.
type duplicateMenuSelection struct {
	allItems bool
	itemIDs  map[uint]struct{}
}

// buildDuplicateMenuSelections determines which menus/items to duplicate based on settings.
func buildDuplicateMenuSelections(sourceMenus []models.Menu, settings *requests.CreateDuplicateVersion) map[uint]*duplicateMenuSelection {
	selections := make(map[uint]*duplicateMenuSelection)

	requestedMenuIDs := make(map[uint]struct{})
	if settings.Menus != nil && len(*settings.Menus) > 0 {
		for i := range *settings.Menus {
			requestedMenuIDs[(*settings.Menus)[i]] = struct{}{}
			selections[(*settings.Menus)[i]] = &duplicateMenuSelection{allItems: true, itemIDs: make(map[uint]struct{})}
		}
	}

	requestedMenuItemIDs := make(map[uint]struct{})
	if settings.MenuItems != nil && len(*settings.MenuItems) > 0 {
		for i := range *settings.MenuItems {
			requestedMenuItemIDs[(*settings.MenuItems)[i]] = struct{}{}
		}
	}

	if len(requestedMenuIDs) == 0 && len(requestedMenuItemIDs) == 0 {
		for i := range sourceMenus {
			selections[sourceMenus[i].ID] = &duplicateMenuSelection{allItems: true, itemIDs: make(map[uint]struct{})}
		}
		return selections
	}

	childRelations := make(map[uint][]models.MenuItemRelation)
	menuParentByChild := make(map[uint]map[uint]sql.Null[uint])

	for i := range sourceMenus {
		menu := sourceMenus[i]
		if _, ok := menuParentByChild[menu.ID]; !ok {
			menuParentByChild[menu.ID] = make(map[uint]sql.Null[uint])
		}

		for j := range menu.MenuItemRelations {
			rel := menu.MenuItemRelations[j]
			childRelations[rel.MenuItemChildID] = append(childRelations[rel.MenuItemChildID], rel)
			menuParentByChild[menu.ID][rel.MenuItemChildID] = rel.MenuItemParentID
		}
	}

	for menuItemID := range requestedMenuItemIDs {
		relations, ok := childRelations[menuItemID]
		if !ok {
			continue
		}

		for i := range relations {
			rel := relations[i]
			selection, ok := selections[rel.MenuID]
			if !ok {
				selection = &duplicateMenuSelection{allItems: false, itemIDs: make(map[uint]struct{})}
				selections[rel.MenuID] = selection
			}

			if selection.allItems {
				continue
			}

			currentID := rel.MenuItemChildID
			for {
				selection.itemIDs[currentID] = struct{}{}

				parent := menuParentByChild[rel.MenuID][currentID]
				if !parent.Valid {
					break
				}

				currentID = parent.V
			}
		}
	}

	return selections
}

// buildCreateMenuItemsFromSource builds a creation DTO tree and a path->oldID map for item remapping.
func buildCreateMenuItemsFromSource(menu *models.Menu, selection *duplicateMenuSelection) (items []requests.CreateMenuItem, oldIDByPath map[string]uint) {
	relationsByParent := make(map[uint][]models.MenuItemRelation)
	oldIDByPath = make(map[string]uint)

	parentByChild := make(map[uint]sql.Null[uint], len(menu.MenuItemRelations))
	relationByChild := make(map[uint]models.MenuItemRelation, len(menu.MenuItemRelations))
	for i := range menu.MenuItemRelations {
		rel := menu.MenuItemRelations[i]
		parentByChild[rel.MenuItemChildID] = rel.MenuItemParentID
		relationByChild[rel.MenuItemChildID] = rel
	}

	for i := range menu.MenuItemRelations {
		rel := menu.MenuItemRelations[i]
		if !selection.allItems {
			if _, ok := selection.itemIDs[rel.MenuItemChildID]; !ok {
				continue
			}
		}

		path := buildMenuItemPath(&rel, relationByChild, parentByChild)
		oldIDByPath[path] = rel.MenuItemChildID

		parentKey := uint(0)
		if rel.MenuItemParentID.Valid {
			parentKey = rel.MenuItemParentID.V
		}
		relationsByParent[parentKey] = append(relationsByParent[parentKey], rel)
	}

	var createItems func(parentID uint) []requests.CreateMenuItem
	createItems = func(parentID uint) []requests.CreateMenuItem {
		relations := relationsByParent[parentID]
		items := make([]requests.CreateMenuItem, 0, len(relations))

		for i := range relations {
			rel := relations[i]
			item := requests.CreateMenuItem{}
			item.SetMenuItemRelation(&rel)

			item.Items = createItems(rel.MenuItemChildID)
			items = append(items, item)
		}

		return items
	}

	items = createItems(0)

	return
}

// buildMenuItemPath creates a deterministic position path for one relation within its menu tree.
func buildMenuItemPath(rel *models.MenuItemRelation, relationByChild map[uint]models.MenuItemRelation, parentByChild map[uint]sql.Null[uint]) string {
	if rel == nil {
		return ""
	}

	segments := make([]string, 0, 4)
	current := *rel

	for {
		segments = append(segments, fmt.Sprintf("%d", current.Position))
		parent := parentByChild[current.MenuItemChildID]
		if !parent.Valid {
			break
		}

		next, ok := relationByChild[parent.V]
		if !ok {
			break
		}
		current = next
	}

	for i, j := 0, len(segments)-1; i < j; i, j = i+1, j-1 {
		segments[i], segments[j] = segments[j], segments[i]
	}

	return strings.Join(segments, "/")
}

// buildMenuItemPathMap creates a position-path lookup for menu item IDs.
func buildMenuItemPathMap(relations []models.MenuItemRelation) map[string]uint {
	pathMap := make(map[string]uint, len(relations))
	parentByChild := make(map[uint]sql.Null[uint], len(relations))
	relationByChild := make(map[uint]models.MenuItemRelation, len(relations))

	for i := range relations {
		rel := relations[i]
		parentByChild[rel.MenuItemChildID] = rel.MenuItemParentID
		relationByChild[rel.MenuItemChildID] = rel
	}

	for i := range relations {
		rel := relations[i]
		path := buildMenuItemPath(&rel, relationByChild, parentByChild)
		pathMap[path] = rel.MenuItemChildID
	}

	return pathMap
}

// duplicatePages clones pages/partials/indexing from old menu items to the new mapped menu items.
func duplicatePages(tx *gorm.DB, menuItemMapping map[uint]uint, locales []string) error {
	for oldMenuItemID, newMenuItemID := range menuItemMapping {
		for i := range locales {
			locale := locales[i]

			sourcePage := &models.Page{}
			result := tx.
				Preload("Indexing").
				Preload("Partials", preloadPagePartialTree).
				Where("menu_item_id = ? AND locale = ?", oldMenuItemID, locale).
				Limit(1).
				Find(sourcePage)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				continue
			}

			targetPage := &models.Page{MenuItemID: newMenuItemID, Locale: sourcePage.Locale, Name: sourcePage.Name}

			if err := tx.Create(targetPage).Error; err != nil {
				return err
			}

			updatePage := requests.UpdatePage{}
			updatePage.SetPage(sourcePage)
			if _, err := UpdatePageWithTx(tx, targetPage, &updatePage); err != nil {
				return err
			}

			for j := range sourcePage.Partials {
				sourcePartial := sourcePage.Partials[j]
				targetPartial := &models.PagePartial{
					MenuItemID: targetPage.MenuItemID,
					Locale:     targetPage.Locale,
					Name:       sourcePartial.Name,
				}

				if err := tx.Create(targetPartial).Error; err != nil {
					return err
				}

				updatePartial := requests.UpdatePagePartial{}
				updatePartial.SetPagePartial(&sourcePartial, targetPartial.ID)
				if _, err := UpdatePagePartialWithTx(tx, targetPartial, &updatePartial); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// duplicateFooter clones footer rows/columns for the selected locales into the target version.
func duplicateFooter(tx *gorm.DB, sourceVersionID, targetVersionID uint, locales []string) error {
	for i := range locales {
		locale := locales[i]
		sourceRows := make([]models.FooterRow, 0)

		if err := preloadFooterTree(tx).
			Where("NOT EXISTS (SELECT 1 FROM footer_row_column_rows frcr WHERE frcr.row_id = footer_rows.id)").
			Order("position asc").
			Find(&sourceRows, "version_id = ? AND locale = ?", sourceVersionID, locale).Error; err != nil {
			return err
		}

		if len(sourceRows) == 0 {
			continue
		}

		dtoRows := make([]requests.UpdateFooterRow, 0, len(sourceRows))
		for j := range sourceRows {
			dtoRow := requests.UpdateFooterRow{}
			dtoRow.SetFooterRow(&sourceRows[j], targetVersionID, locale)
			dtoRows = append(dtoRows, dtoRow)
		}

		existingRows := make([]models.FooterRow, 0)
		if _, err := UpdateFooterWithTx(tx, targetVersionID, locale, &existingRows, &requests.UpdateFooter{Rows: dtoRows}); err != nil {
			return err
		}
	}

	return nil
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
