package services

import (
	"api-page/main/src/cache"
	"api-page/main/src/database"
	"api-page/main/src/dto/requests"
	"api-page/main/src/dto/responses"
	"api-page/main/src/enums"
	"api-page/main/src/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/gofiber/fiber/v2"
	"github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IsMenuNameAvailable method to check if a menu is available.
func IsMenuNameAvailable(versionID uint, name string, ignore *string) (bool, error) {
	query := database.Pg.Limit(1)
	var result *gorm.DB
	if ignore != nil {
		result = query.Find(&models.Menu{}, "version_id = ? AND name = ? AND name != ?", versionID, name, ignore)
	} else {
		result = query.Find(&models.Menu{}, "version_id = ? AND name = ?", versionID, name)
	}

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected == 0, nil
}

// IsMenuDeleted method to check if a menu is deleted.
func IsMenuDeleted(menuID uint) (bool, error) {
	if result := database.Pg.Unscoped().Limit(1).Find(&models.Menu{}, "id = ? AND deleted_at IS NOT NULL", menuID); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// IsMenuItemWithAppName method to check if a menu item belongs to a version with the given app name.
func IsMenuItemWithAppName(menuItemID uint, appName string) (bool, error) {
	var count int64

	if result := database.Pg.Model(&models.MenuItem{}).
		Joins("JOIN menu_item_relations ON menu_item_relations.menu_item_child_id = menu_items.id").
		Joins("JOIN menus ON menus.id = menu_item_relations.menu_id").
		Joins("JOIN versions ON versions.id = menus.version_id").
		Where("menu_items.id = ? AND versions.app_name = ?", menuItemID, appName).
		Count(&count); result.Error != nil {
		return false, result.Error
	}

	return count > 0, nil
}

// GetMenus method to get paginated menus.
func GetMenus(c *fiber.Ctx) (*pagination.Model, error) {
	menus := make([]models.Menu, 0)
	values := c.Request().URI().QueryArgs()
	allowedColumns := map[string]bool{
		"id":         true,
		"version_id": true,
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
	dbResult := database.Pg.Scopes(queryFunc, sortFunc).
		Limit(limit).
		Offset(offset)

	total := int64(0)
	dbCount := database.Pg.Scopes(queryFunc).
		Model(&models.Menu{})

	if result := dbResult.Find(&menus); result.Error != nil {
		return nil, result.Error
	}

	dbCount.Count(&total)
	pageCount := pagination.Count(int(total), limit)

	paginatedMenus := make([]responses.PaginatedMenu, 0)
	for i := range menus {
		paginatedMenu := responses.PaginatedMenu{}
		paginatedMenu.SetPaginatedMenu(&menus[i])
		paginatedMenus = append(paginatedMenus, paginatedMenu)
	}

	paginationModel := pagination.CreatePaginationModel(limit, page, pageCount, int(total), paginatedMenus)

	return &paginationModel, nil
}

// GetMenuLookup method to get a lookup of menus.
func GetMenuLookup(versionID uint, name *string) (*[]models.Menu, error) {
	menus := make([]models.Menu, 0)

	if inCache, err := isMenusLookupInCache(versionID); err != nil {
		return nil, err
	} else if inCache {
		if cacheMenus, err := getMenusLookupFromCache(versionID); err != nil {
			return nil, err
		} else if cacheMenus != nil && len(*cacheMenus) > 0 {
			menus = *cacheMenus
		}
	}

	if len(menus) == 0 {
		query := database.Pg.Model(&models.Menu{}).
			Select("id", "name")

		if result := query.Find(&menus, "version_id = ?", versionID); result.Error != nil {
			return nil, result.Error
		}

		_ = setMenusLookupToCache(versionID, &menus)
	}

	// If a name filter is provided, perform case-insensitive substring match on the list.
	if name != nil {
		target := strings.TrimSpace(*name)
		if target != "" {
			lowerTarget := strings.ToLower(target)
			filtered := make([]models.Menu, 0, len(menus))
			for i := range menus {
				if strings.Contains(strings.ToLower(menus[i].Name), lowerTarget) {
					filtered = append(filtered, menus[i])
				}
			}
			menus = filtered
		}
	}

	return &menus, nil
}

// GetMenusByVersionID method to get menus by version ID and locale.
func GetMenusByVersionID(versionID uint, locale string) (*[]models.Menu, error) {
	menus := make([]models.Menu, 0)

	if inCache, err := isVersionMenusInCache(versionID); err != nil {
		return nil, err
	} else if inCache {
		if cacheMenus, err := getVersionMenusFromCache(versionID, locale); err != nil {
			return nil, err
		} else if cacheMenus != nil && len(*cacheMenus) > 0 {
			menus = *cacheMenus
		}
	}

	if len(menus) == 0 {
		if result := database.Pg.
			Preload("MenuItemRelations", func(db *gorm.DB) *gorm.DB {
				return db.Preload("MenuItemChild", func(db2 *gorm.DB) *gorm.DB {
					return db2.Preload("Pages", func(db3 *gorm.DB) *gorm.DB {
						return db3.Where("locale = ? AND enabled_at IS NOT NULL", locale)
					}).Where("enabled_at IS NOT NULL")
				}).
					Joins("JOIN menu_items mi ON mi.id = menu_item_relations.menu_item_child_id").
					Joins("JOIN pages p ON p.menu_item_id = mi.id AND p.locale = ? AND p.enabled_at IS NOT NULL", locale).
					Order("menu_item_parent_id NULLS FIRST").
					Order("position ASC")
			}).
			Find(&menus, "version_id = ?", versionID); result.Error != nil {
			return nil, result.Error
		}

		_ = setVersionMenusToCache(versionID, locale, &menus)
	}

	return &menus, nil
}

// GetMenuByID method to get a menu by ID.
func GetMenuByID(menuID uint) (*models.Menu, error) {
	menu := &models.Menu{}

	if result := database.Pg.
		Preload("MenuItemRelations", func(db *gorm.DB) *gorm.DB {
			return db.Preload("MenuItemParent", func(db2 *gorm.DB) *gorm.DB { return db2.Preload("Indexing") }).
				Preload("MenuItemChild", func(db2 *gorm.DB) *gorm.DB { return db2.Preload("Indexing") }).
				Order("menu_item_parent_id NULLS FIRST").
				Order("position ASC")
		}).
		Find(menu, "id = ?", menuID); result.Error != nil {
		return nil, result.Error
	}

	return menu, nil
}

// GetVersionIDByMenuItemID method to get version ID by menu item ID.
func GetVersionIDByMenuItemID(menuItemID uint) (uint, error) {
	var versionID uint

	if result := database.Pg.Model(&models.MenuItem{}).Where("id = ?", menuItemID).Pluck("version_id", &versionID); result.Error != nil {
		return 0, result.Error
	}

	return versionID, nil
}

// CreateMenu method to create a menu.
func CreateMenu(menu *requests.CreateMenu) (*models.Menu, error) {
	m := &models.Menu{VersionID: menu.VersionID, Name: menu.Name}

	result := &models.Menu{}
	if err := database.Pg.Transaction(func(tx *gorm.DB) error {
		if err := tx.FirstOrCreate(&result, m).Error; err != nil {
			return err
		}

		// Ensure slice is initialized.
		result.MenuItemRelations = make([]models.MenuItemRelation, 0)

		// Persist hierarchical menu items and their relations while populating result.
		for i := range menu.Items {
			if _, err := createMenuItemHierarchy(tx, result, nil, &menu.Items[i]); err != nil {
				return err
			}
		}

		// Sort relations grouped by parent and by position.
		sortMenuItemRelations(result.MenuItemRelations)

		return nil
	}); err != nil {
		return nil, err
	}

	_ = deleteMenusLookupFromCache(m.VersionID)
	_ = deleteAllVersionMenusFromCache(m.VersionID)

	return result, nil
}

// UpdateMenu method to update a menu.
func UpdateMenu(oldMenu *models.Menu, menu *requests.UpdateMenu) (*models.Menu, error) {
	// Persist updates in a transaction to keep tree consistent.
	if err := database.Pg.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&oldMenu).
			Clauses(clause.Returning{Columns: []clause.Column{{Name: "name"}, {Name: "updated_at"}}}).
			Updates(models.Menu{Name: menu.Name}).Error; err != nil {
			return err
		}

		// Reset in-memory relations; we'll rebuild from DTO.
		oldMenu.MenuItemRelations = make([]models.MenuItemRelation, 0)

		// Reconcile top-level items (parent is nil).
		if err := reconcileMenuItems(tx, oldMenu, nil, menu.Items); err != nil {
			return err
		}

		// Sort relations grouped by parent and by position.
		sortMenuItemRelations(oldMenu.MenuItemRelations)

		return nil
	}); err != nil {
		return nil, err
	}

	_ = deleteMenusLookupFromCache(oldMenu.VersionID)
	_ = deleteAllVersionMenusFromCache(oldMenu.VersionID)

	return oldMenu, nil
}

// DeleteMenu method to delete a menu.
func DeleteMenu(versionID, menuID uint) error {
	err := database.Pg.Delete(&models.Menu{}, menuID).Error
	if err == nil {
		_ = deleteMenusLookupFromCache(versionID)
		_ = deleteAllVersionMenusFromCache(versionID)
	}

	return err
}

// RestoreMenu method to restore a deleted menu.
func RestoreMenu(menuID uint) error {
	err := database.Pg.Unscoped().Model(&models.Menu{}).Where("id = ?", menuID).Update("deleted_at", nil).Error
	if err == nil {
		var versionID uint

		if result := database.Pg.Unscoped().Model(&models.Menu{}).Where("id = ?", menuID).Pluck("version_id", &versionID); result.Error != nil {
			return result.Error
		}

		_ = deleteMenusLookupFromCache(versionID)
		_ = deleteAllVersionMenusFromCache(versionID)
	}

	return err
}

// getMenusLookupCacheKey gets the key for the cache.
func getMenusLookupCacheKey(versionID uint) string {
	return fmt.Sprintf("menus:lookup:%d", versionID)
}

// isMenusLookupInCache checks if the menus exists in the cache.
func isMenusLookupInCache(versionID uint) (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(getMenusLookupCacheKey(versionID)).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// getMenusLookupFromCache gets the menus from the cache.
func getMenusLookupFromCache(versionID uint) (*[]models.Menu, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(getMenusLookupCacheKey(versionID)).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var menus []models.Menu
	if err := json.Unmarshal([]byte(value), &menus); err != nil {
		return nil, err
	}

	return &menus, nil
}

// setMenusLookupToCache sets the menus to the cache.
func setMenusLookupToCache(versionID uint, menus *[]models.Menu) error {
	value, err := json.Marshal(menus)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(getMenusLookupCacheKey(versionID)).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deleteMenusLookupFromCache deletes existing menus from the cache.
func deleteMenusLookupFromCache(versionID uint) error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(getMenusLookupCacheKey(versionID)).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// getVersionMenusCacheKey gets the key for the cache.
func getVersionMenusCacheKey(versionID uint) string {
	return fmt.Sprintf("menus:version:%d:locales", versionID)
}

// isVersionMenusInCache checks if the menus of a version exists in the cache.
func isVersionMenusInCache(versionID uint) (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(getVersionMenusCacheKey(versionID)).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// getAllVersionMenusFromCache gets the menus in a version from the cache.
func getAllVersionMenusFromCache(versionID uint) (*map[string][]models.Menu, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(getVersionMenusCacheKey(versionID)).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var versionMenus map[string][]models.Menu
	if err := json.Unmarshal([]byte(value), &versionMenus); err != nil {
		return nil, err
	}

	return &versionMenus, nil
}

// getVersionMenusFromCache gets the menus in a version with a locale from the cache.
func getVersionMenusFromCache(versionID uint, locale string) (*[]models.Menu, error) {
	var versionMenus *map[string][]models.Menu

	if inCache, err := isVersionMenusInCache(versionID); err != nil {
		return nil, err
	} else if inCache {
		if versionMenus, err = getAllVersionMenusFromCache(versionID); err != nil {
			return nil, err
		}
	}

	if versionMenus == nil {
		return nil, nil
	}

	if menus, ok := (*versionMenus)[locale]; ok {
		return &menus, nil
	}

	return nil, nil
}

// setVersionMenusToCache sets the menus of a version to the cache.
func setVersionMenusToCache(versionID uint, locale string, menus *[]models.Menu) error {
	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	var versionMenus *map[string][]models.Menu

	if inCache, err := isVersionMenusInCache(versionID); err != nil {
		return err
	} else if inCache {
		if versionMenus, err = getAllVersionMenusFromCache(versionID); err != nil {
			return err
		}
	}

	if versionMenus == nil {
		versionMenus = &map[string][]models.Menu{}
		(*versionMenus)[locale] = *menus
	} else {
		(*versionMenus)[locale] = *menus
	}

	value, err := json.Marshal(versionMenus)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(getVersionMenusCacheKey(versionID)).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deleteVersionMenusFromCache deletes existing menus in a version from the cache.
func deleteVersionMenusFromCache(versionID uint, locale string) error {
	var versionMenus *map[string][]models.Menu

	if inCache, err := isVersionMenusInCache(versionID); err != nil {
		return err
	} else if inCache {
		if versionMenus, err = getAllVersionMenusFromCache(versionID); err != nil {
			return err
		}
	}

	if versionMenus == nil {
		return nil
	}

	delete(*versionMenus, locale)

	if len(*versionMenus) > 0 {
		expiration := os.Getenv("VALKEY_EXPIRATION")
		duration, err := time.ParseDuration(expiration)
		if err != nil {
			return err
		}

		value, err := json.Marshal(versionMenus)
		if err != nil {
			return err
		}

		result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(getVersionMenusCacheKey(versionID)).Value(valkey.BinaryString(value)).Ex(duration).Build())
		if result.Error() != nil {
			return result.Error()
		}

		return nil
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(getVersionMenusCacheKey(versionID)).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deleteAllVersionMenusFromCache deletes existing menus in a version for all languages from the cache.
func deleteAllVersionMenusFromCache(versionID uint) error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(getVersionMenusCacheKey(versionID)).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// sortMenuItemRelations sorts the relations grouped by parent (NULL first, then by parent ID) and within each group by Position.
func sortMenuItemRelations(relations []models.MenuItemRelation) {
	sort.SliceStable(relations, func(i, j int) bool {
		pi := relations[i].MenuItemParentID
		pj := relations[j].MenuItemParentID

		// Group by parent: nulls first, then ascending by parent ID.
		if pi.Valid != pj.Valid {
			return !pi.Valid && pj.Valid
		}
		if pi.Valid && pj.Valid {
			if pi.V != pj.V {
				return pi.V < pj.V
			}
		}
		// Same parent group: sort by position ascending.
		if relations[i].Position != relations[j].Position {
			return relations[i].Position < relations[j].Position
		}
		// Fallback: by child ID to stabilize.
		return relations[i].MenuItemChildID < relations[j].MenuItemChildID
	})
}

// reconcileMenuItems upserts items/relations/indexing under a given parent, and deletes removed children.
func reconcileMenuItems(tx *gorm.DB, menu *models.Menu, parent *models.MenuItem, items []requests.UpdateMenuItem) error {
	processed := make(map[uint]struct{})

	for i := range items {
		dto := &items[i]

		mi, err := upsertMenuItem(tx, menu, dto)
		if err != nil {
			return err
		}

		// Ensure the relation reflects desired parent and position.
		desiredParentID := sql.Null[uint]{}
		if parent != nil {
			desiredParentID = sql.Null[uint]{V: parent.ID, Valid: true}
		}

		// Find any existing relation for this child in this menu.
		existingRel := &models.MenuItemRelation{}
		if err := tx.Where("menu_id = ? AND menu_item_child_id = ?", menu.ID, mi.ID).Find(&existingRel).Error; err != nil {
			return err
		}

		// Remove relations that don't match desired parent and position.
		if existingRel.MenuID != 0 {
			matchesParent := existingRel.MenuItemParentID.Valid == desiredParentID.Valid && (!existingRel.MenuItemParentID.Valid || existingRel.MenuItemParentID.V == desiredParentID.V)
			if !(matchesParent && existingRel.Position == *dto.Position) {
				if err := tx.Delete(&existingRel).Error; err != nil {
					return err
				}
			}
		}

		// Ensure target relation exists.
		rel := &models.MenuItemRelation{MenuID: menu.ID, MenuItemChildID: mi.ID, Position: *dto.Position}
		if desiredParentID.Valid {
			rel.MenuItemParentID = desiredParentID
		}
		if err := tx.FirstOrCreate(rel, rel).Error; err != nil {
			return err
		}

		// Populate and append to in-memory menu
		rel.MenuItemChild = *mi
		if parent != nil {
			rel.MenuItemParent = *parent
		}
		menu.MenuItemRelations = append(menu.MenuItemRelations, *rel)

		// Recurse into children.
		if err := reconcileMenuItems(tx, menu, mi, dto.Items); err != nil {
			return err
		}

		processed[mi.ID] = struct{}{}
	}

	// Delete children that are no longer present under this parent.
	toDelete := make([]models.MenuItemRelation, 0)
	base := tx.Where("menu_id = ?", menu.ID)
	if parent == nil {
		base = base.Where("menu_item_parent_id IS NULL")
	} else {
		base = base.Where("menu_item_parent_id = ?", parent.ID)
	}
	if err := base.Find(&toDelete).Error; err != nil {
		return err
	}
	for i := range toDelete {
		childID := toDelete[i].MenuItemChildID
		if _, ok := processed[childID]; !ok {
			// Soft delete the MenuItem (cascades will remove relations and indexing via constraints).
			if err := tx.Delete(&models.MenuItem{}, childID).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

// upsertMenuItem creates or updates a MenuItem and synchronizes its indexing.
func upsertMenuItem(tx *gorm.DB, menu *models.Menu, dto *requests.UpdateMenuItem) (*models.MenuItem, error) {
	mi := &models.MenuItem{Indexing: make([]models.MenuItemIndexing, 0)}
	isNew := dto.ID == nil
	if !isNew {
		if err := tx.First(mi, "id = ?", *dto.ID).Error; err != nil {
			return nil, err
		}

		mi.Name = dto.Name
		if dto.Icon != nil {
			mi.Icon = sql.NullString{String: *dto.Icon, Valid: true}
		} else {
			mi.Icon = sql.NullString{Valid: false}
		}
		if dto.EnabledAt != nil {
			mi.EnabledAt = sql.NullTime{Time: *dto.EnabledAt, Valid: true}
		} else {
			mi.EnabledAt = sql.NullTime{Valid: false}
		}
		if err := tx.Save(mi).Error; err != nil {
			return nil, err
		}
		if err := deletePagesFromCacheByMenuItemID(mi.ID); err != nil {
			return nil, err
		}
	} else {
		mi = &models.MenuItem{VersionID: menu.VersionID, Name: dto.Name}
		if dto.Icon != nil {
			mi.Icon = sql.NullString{String: *dto.Icon, Valid: true}
		}
		if dto.EnabledAt != nil {
			mi.EnabledAt = sql.NullTime{Time: *dto.EnabledAt, Valid: true}
		}
		if err := tx.FirstOrCreate(mi, mi).Error; err != nil {
			return nil, err
		}
	}

	// Sync indexing: upsert present options and remove missing ones.
	existing := make([]models.MenuItemIndexing, 0)
	if err := tx.Where("menu_item_id = ?", mi.ID).Find(&existing).Error; err != nil {
		return nil, err
	}

	desired := make(map[enums.Indexing]*string, len(dto.Indexing))
	for i := range dto.Indexing {
		opt := enums.Indexing(dto.Indexing[i].Option)
		desired[opt] = dto.Indexing[i].Value
		mii := &models.MenuItemIndexing{MenuItemID: mi.ID, Option: opt}
		if dto.Indexing[i].Value != nil {
			mii.Value = sql.NullString{String: *dto.Indexing[i].Value, Valid: true}
		} else {
			mii.Value = sql.NullString{Valid: false}
		}
		if err := tx.FirstOrCreate(mii, mii).Error; err != nil {
			return nil, err
		}
		// Ensure value is updated when it already existed.
		if err := tx.Model(&models.MenuItemIndexing{}).
			Where("menu_item_id = ? AND option = ?", mi.ID, opt).
			Updates(map[string]interface{}{"value": mii.Value}).Error; err != nil {
			return nil, err
		}

		mi.Indexing = append(mi.Indexing, *mii)
	}

	for i := range existing {
		if _, ok := desired[existing[i].Option]; !ok {
			if err := tx.Delete(&existing[i]).Error; err != nil {
				return nil, err
			}
		}
	}

	return mi, nil
}

// createMenuItemHierarchy creates/gets a MenuItem, its indexing, the relation entry, and recurses children.
// It also appends a flattened MenuItemRelation (with populated parent/child) to the provided result menu.
func createMenuItemHierarchy(tx *gorm.DB, menu *models.Menu, parent *models.MenuItem, item *requests.CreateMenuItem) (uint, error) {
	mi := &models.MenuItem{
		VersionID: menu.VersionID,
		Name:      item.Name,
		Indexing:  make([]models.MenuItemIndexing, len(item.Indexing)),
	}
	if item.Icon != nil {
		mi.Icon = sql.NullString{String: *item.Icon, Valid: true}
	}
	if item.EnabledAt != nil {
		mi.EnabledAt = sql.NullTime{Time: *item.EnabledAt, Valid: true}
	}
	if err := tx.FirstOrCreate(mi, mi).Error; err != nil {
		return 0, err
	}

	for i := range item.Indexing {
		idx := item.Indexing[i]
		mii := &models.MenuItemIndexing{MenuItemID: mi.ID, Option: enums.Indexing(idx.Option)}
		if idx.Value != nil {
			mii.Value = sql.NullString{String: *idx.Value, Valid: true}
		}
		if err := tx.FirstOrCreate(mii, mii).Error; err != nil {
			return 0, err
		}
		mi.Indexing[i] = *mii
	}

	rel := &models.MenuItemRelation{
		MenuID:          menu.ID,
		MenuItemChildID: mi.ID,
		Position:        *item.Position,
	}
	if parent != nil {
		rel.MenuItemParentID = sql.Null[uint]{V: parent.ID, Valid: true}
	}
	if err := tx.FirstOrCreate(rel, rel).Error; err != nil {
		return 0, err
	}

	rel.MenuItemChild = *mi
	if parent != nil {
		rel.MenuItemParent = *parent
	}
	menu.MenuItemRelations = append(menu.MenuItemRelations, *rel)

	// Recurse for children.
	for i := range item.Items {
		if _, err := createMenuItemHierarchy(tx, menu, mi, &item.Items[i]); err != nil {
			return 0, err
		}
	}

	return mi.ID, nil
}
