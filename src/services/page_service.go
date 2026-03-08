package services

import (
	"api-page/main/src/cache"
	"api-page/main/src/database"
	"api-page/main/src/dto/requests"
	"api-page/main/src/enums"
	"api-page/main/src/models"
	"api-page/main/src/utils"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const MaxPagePartialTreeDepth = 4

// IsPagePartialNameAvailable method to check if a name of a partial is available.
func IsPagePartialNameAvailable(menuItemID uint, locale, name string, ignore *string) (bool, error) {
	query := database.Pg.Limit(1)
	var result *gorm.DB
	if ignore != nil {
		result = query.Find(&models.PagePartial{}, "menu_item_id = ? AND locale = ? AND name = ? AND name != ?", menuItemID, locale, name, ignore)
	} else {
		result = query.Find(&models.PagePartial{}, "menu_item_id = ? AND locale = ? AND name = ?", menuItemID, locale, name)
	}

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected == 0, nil
}

// IsPageDeleted method to check if a page or any of its related records are deleted.
// It returns true if the Page (by menu_item_id & locale), the MenuItem (by id), or any related Menu
// (via menu_item_relations -> menus) has a non-null deleted_at.
func IsPageDeleted(menuItemID uint, locale string) (bool, error) {
	// Single SQL using EXISTS for performance and clarity.
	// Checks:
	// 1) Page deleted for given menu_item_id + locale.
	// 2) MenuItem deleted for given id.
	// 3) Any Menu deleted that is related to the given menu item through menu_item_relations.
	const q = `
		SELECT (
			EXISTS(
				SELECT 1
				FROM pages p
				WHERE p.menu_item_id = ? AND p.locale = ? AND p.deleted_at IS NOT NULL
			)
		) OR (
			EXISTS(
				SELECT 1
				FROM menu_items mi
				WHERE mi.id = ? AND mi.deleted_at IS NOT NULL
			)
		) OR (
			EXISTS(
				SELECT 1
				FROM menu_item_relations mir
				JOIN menus m ON m.id = mir.menu_id
				WHERE mir.menu_item_child_id = ? AND m.deleted_at IS NOT NULL
			)
		) AS any_deleted;
	`

	var anyDeleted bool
	row := database.Pg.Raw(q, menuItemID, locale, menuItemID, menuItemID).Row()
	if err := row.Scan(&anyDeleted); err != nil {
		return false, err
	}

	return anyDeleted, nil
}

// IsPagePartialDeleted method to check if a page partial is deleted by its ID.
func IsPagePartialDeleted(partialID uint) (bool, error) {
	partial := &models.PagePartial{}
	if result := database.Pg.Unscoped().Find(partial, "id = ?", partialID); result.Error != nil {
		return false, result.Error
	}

	return partial.DeletedAt.Valid, nil
}

// IsLastEnabledPagePartial method to check if the page partial is the last enabled partial for the given menu item and locale.
func IsLastEnabledPagePartial(menuItemID uint, locale string) (bool, error) {
	var count int64

	if result := database.Pg.Model(&models.PagePartial{}).
		Where("menu_item_id = ? AND locale = ?", menuItemID, locale).
		Count(&count); result.Error != nil {
		return false, result.Error
	}

	return count == 1, nil
}

// GetPage retrieves a Page by MenuItemID and Locale.
func GetPage(menuItemID uint, locale string) (*models.Page, error) {
	page := &models.Page{}

	if isPageDeleted, err := IsPageDeleted(menuItemID, locale); err != nil {
		return nil, err
	} else if isPageDeleted {
		// If the page is deleted, we should not retrieve it.
		return page, nil
	}

	if result := database.Pg.Find(page, "menu_item_id = ? AND locale = ?", menuItemID, locale); result.Error != nil {
		return nil, result.Error
	}

	return page, nil
}

// GetPublishedPage retrieves a Page by MenuItemID and Locale only if it is enabled (EnabledAt is not null) and not deleted.
func GetPublishedPage(menuItemID uint, locale string) (*models.Page, error) {
	page := &models.Page{}

	if inCache, err := isPageInCache(menuItemID, locale); err != nil {
		return nil, err
	} else if inCache {
		if cachePage, err := getPageFromCache(menuItemID, locale); err != nil {
			return nil, err
		} else if cachePage != nil {
			page = cachePage
		}
	}

	if page.MenuItemID == 0 {
		if isPageDeleted, err := IsPageDeleted(menuItemID, locale); err != nil {
			return nil, err
		} else if isPageDeleted {
			// If the page is deleted, we should not retrieve it.
			return page, nil
		}

		if result := database.Pg.
			Preload("Indexing").
			Preload("Partials", preloadPagePartialTree).
			Find(page, "menu_item_id = ? AND locale = ? AND enabled_at IS NOT NULL", menuItemID, locale); result.Error != nil {
			return nil, result.Error
		}

		if page.EnabledAt.Valid == false {
			// If the page is not enabled, we should not retrieve it.
			return &models.Page{}, nil
		}

		if len(page.Indexing) == 0 {
			menuIndexing := make([]models.MenuItemIndexing, 0)

			if result := database.Pg.Find(&menuIndexing, "menu_item_id = ?", menuItemID); result.Error != nil {
				return nil, result.Error
			}

			for i := range menuIndexing {
				pageIndexing := models.PageIndexing{
					MenuItemID: menuIndexing[i].MenuItemID,
					Locale:     locale,
					Option:     menuIndexing[i].Option,
					Value:      menuIndexing[i].Value,
				}
				page.Indexing = append(page.Indexing, pageIndexing)
			}
		}

		_ = setPageToCache(menuItemID, locale, page)
	}

	return page, nil
}

// GetOrCreatePage retrieves a Page by MenuItemID and Locale. If it doesn't exist, it creates a new one.
func GetOrCreatePage(menuItemID uint, locale string) (*models.Page, error) {
	page := &models.Page{}

	if isPageDeleted, err := IsPageDeleted(menuItemID, locale); err != nil {
		return nil, err
	} else if isPageDeleted {
		// If the page is deleted, we should not retrieve or create it.
		return page, nil
	}

	page.MenuItemID = menuItemID
	page.Locale = locale
	page.Name = ""

	result := database.Pg.
		Preload("Indexing").
		Preload("Partials", preloadPagePartialTree).
		FirstOrCreate(page, page)
	if result.Error != nil {
		return nil, result.Error
	}

	if page.Name == "" && len(page.Partials) == 0 {
		partial := &models.PagePartial{
			MenuItemID: page.MenuItemID,
			Locale:     page.Locale,
			Name:       "Default",
		}

		if err := database.Pg.FirstOrCreate(partial, partial).Error; err != nil {
			return nil, err
		}

		page.Partials = append(page.Partials, *partial)
	}

	// Invalidate cache for version menus related to this page.
	versionID, err := GetVersionIDByMenuItemID(menuItemID)
	if err != nil {
		return nil, err
	}
	_ = deleteVersionMenusFromCache(versionID, locale)

	return page, nil
}

// GetPartialByID retrieves a PagePartial by its ID, including its associated rows and columns.
func GetPartialByID(partialID uint) (*models.PagePartial, error) {
	partial := &models.PagePartial{}

	if result := preloadPagePartialTree(database.Pg).
		Find(partial, "id = ?", partialID); result.Error != nil {
		return nil, result.Error
	}

	return partial, nil
}

// CreatePagePartial creates a new PagePartial for the given Page using data from the CreatePagePartial request.
func CreatePagePartial(page *models.Page, request *requests.CreatePagePartial) (*models.PagePartial, error) {
	partial := &models.PagePartial{
		MenuItemID: page.MenuItemID,
		Locale:     page.Locale,
		Name:       request.Name,
	}

	if err := database.Pg.FirstOrCreate(partial, partial).Error; err != nil {
		return nil, err
	}

	_ = deletePageFromCache(page.MenuItemID, page.Locale)

	return partial, nil
}

// UpdatePage updates the given Page with data from the UpdatePage request.
func UpdatePage(page *models.Page, request *requests.UpdatePage) (*models.Page, error) {
	if isPageDeleted, err := IsPageDeleted(page.MenuItemID, page.Locale); err != nil {
		return nil, err
	} else if isPageDeleted {
		// If the page is deleted, we should not update it.
		return nil, gorm.ErrRecordNotFound
	}

	page.Name = request.Name
	page.NewTabEnabled = request.NewTabEnabled
	page.UrlEnabled = request.UrlEnabled
	page.Plugin = utils.NewNullString(request.Plugin)
	page.MetaTitle = utils.NewNullString(request.MetaTitle)
	page.MetaDescription = utils.NewNullString(request.MetaDescription)
	page.Hashtag = utils.NewNullString(request.Hashtag)
	page.Url = utils.NewNullString(request.Url)
	page.EnabledAt = utils.NewNullTime(request.EnabledAt)
	page.Indexing = make([]models.PageIndexing, 0)

	if err := database.Pg.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&page).
			Clauses(clause.Returning{Columns: []clause.Column{{Name: "updated_at"}}}).
			Updates(page).Error; err != nil {
			tx.Rollback()
			return err
		}

		existing := make([]models.PageIndexing, 0)
		if err := tx.Where("menu_item_id = ? AND locale = ?", page.MenuItemID, page.Locale).Find(&existing).Error; err != nil {
			tx.Rollback()
			return err
		}

		desired := make(map[enums.Indexing]*string, len(request.Indexing))
		for i := range request.Indexing {
			opt := enums.Indexing(request.Indexing[i].Option)
			desired[opt] = request.Indexing[i].Value
			pi := &models.PageIndexing{MenuItemID: page.MenuItemID, Locale: page.Locale, Option: opt}
			pi.Value = utils.NewNullString(request.Indexing[i].Value)
			if err := tx.FirstOrCreate(pi, pi).Error; err != nil {
				tx.Rollback()
				return err
			}
			// Ensure value is updated when it already existed.
			if err := tx.Model(&models.PageIndexing{}).
				Where("menu_item_id = ? AND locale = ? AND option = ?", page.MenuItemID, page.Locale, opt).
				Updates(map[string]interface{}{"value": pi.Value}).Error; err != nil {
				tx.Rollback()
				return err
			}

			page.Indexing = append(page.Indexing, *pi)
		}

		for i := range existing {
			if _, ok := desired[existing[i].Option]; !ok {
				if err := tx.Delete(&existing[i]).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// Invalidate cache for version menus related to this page.
	versionID, err := GetVersionIDByMenuItemID(page.MenuItemID)
	if err != nil {
		return nil, err
	}
	_ = deleteVersionMenusFromCache(versionID, page.Locale)
	_ = deletePageFromCache(page.MenuItemID, page.Locale)

	return page, nil
}

// UpdatePagePartial updates the given PagePartial and its associated rows and columns
// based on the data provided in the UpdatePagePartial request.
func UpdatePagePartial(menuItemID uint, locale string, partial *models.PagePartial, dtoPartial *requests.UpdatePagePartial) (*models.PagePartial, error) {
	partial.Name = dtoPartial.Name

	if err := database.Pg.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(partial).
			Clauses(clause.Returning{Columns: []clause.Column{{Name: "updated_at"}}}).
			Updates(partial).Error; err != nil {
			return err
		}

		existingRows := make([]models.PagePartialRow, len(partial.Rows))
		copy(existingRows, partial.Rows)

		rows, err := syncPagePartialRows(tx, partial.ID, nil, existingRows, dtoPartial.Rows, 1)
		if err != nil {
			return err
		}

		partial.Rows = rows
		return nil
	}); err != nil {
		return nil, err
	}

	_ = deletePageFromCache(menuItemID, locale)

	return partial, nil
}

// syncPagePartialRows synchronizes the PagePartialRows for a given PagePartial based on the provided DTO rows.
func syncPagePartialRows(tx *gorm.DB, partialID uint, parentColumnID *uint, existingRows []models.PagePartialRow, dtoRows []requests.UpdatePagePartialRow, depth int) ([]models.PagePartialRow, error) {
	if depth > MaxPagePartialTreeDepth {
		return nil, fmt.Errorf("page partial row depth exceeded max depth of %d", MaxPagePartialTreeDepth)
	}

	existingByID := make(map[uint]*models.PagePartialRow, len(existingRows))
	for i := range existingRows {
		existingByID[existingRows[i].ID] = &existingRows[i]
	}

	desiredRows := make([]models.PagePartialRow, 0, len(dtoRows))
	desiredIDs := make(map[uint]bool, len(dtoRows))

	for i := range dtoRows {
		dtoRow := dtoRows[i]
		var row *models.PagePartialRow
		if dtoRow.ID != nil {
			row = existingByID[*dtoRow.ID]
		}
		if row == nil {
			row = &models.PagePartialRow{}
		}

		row.PagePartialID = partialID
		row.Position = utils.UintOrZero(dtoRow.Position)
		row.NoGutters = dtoRow.NoGutters
		row.Dense = dtoRow.Dense
		row.Hashtag = utils.NewNullString(dtoRow.Hashtag)
		row.Align = utils.NewNullString(dtoRow.Align)
		row.AlignXxl = utils.NewNullString(dtoRow.AlignXxl)
		row.AlignXl = utils.NewNullString(dtoRow.AlignXl)
		row.AlignLg = utils.NewNullString(dtoRow.AlignLg)
		row.AlignMd = utils.NewNullString(dtoRow.AlignMd)
		row.AlignSm = utils.NewNullString(dtoRow.AlignSm)
		row.AlignContent = utils.NewNullString(dtoRow.AlignContent)
		row.AlignContentXxl = utils.NewNullString(dtoRow.AlignContentXxl)
		row.AlignContentXl = utils.NewNullString(dtoRow.AlignContentXl)
		row.AlignContentLg = utils.NewNullString(dtoRow.AlignContentLg)
		row.AlignContentMd = utils.NewNullString(dtoRow.AlignContentMd)
		row.AlignContentSm = utils.NewNullString(dtoRow.AlignContentSm)
		row.Justify = utils.NewNullString(dtoRow.Justify)
		row.JustifyXxl = utils.NewNullString(dtoRow.JustifyXxl)
		row.JustifyXl = utils.NewNullString(dtoRow.JustifyXl)
		row.JustifyLg = utils.NewNullString(dtoRow.JustifyLg)
		row.JustifyMd = utils.NewNullString(dtoRow.JustifyMd)
		row.JustifySm = utils.NewNullString(dtoRow.JustifySm)

		if row.ID != 0 {
			if err := tx.Model(row).
				Clauses(clause.Returning{Columns: []clause.Column{{Name: "updated_at"}}}).
				Updates(row).Error; err != nil {
				return nil, err
			}
		} else {
			if err := tx.Create(row).Error; err != nil {
				return nil, err
			}
		}

		if parentColumnID != nil {
			if err := ensureRowColumnRelation(tx, *parentColumnID, row.ID); err != nil {
				return nil, err
			}
		}

		existingColumns := make([]models.PagePartialRowColumn, 0)
		if dtoRow.ID != nil && len(row.Columns) == 0 {
			if err := tx.Where("page_partial_row_id = ?", row.ID).Find(&existingColumns).Error; err != nil {
				return nil, err
			}
		} else {
			existingColumns = row.Columns
		}

		columns, err := syncPagePartialColumns(tx, partialID, row.ID, existingColumns, dtoRow.Columns, depth)
		if err != nil {
			return nil, err
		}

		row.Columns = columns
		desiredRows = append(desiredRows, *row)
		desiredIDs[row.ID] = true
	}

	for i := range existingRows {
		if _, ok := desiredIDs[existingRows[i].ID]; !ok {
			if err := tx.Delete(&existingRows[i]).Error; err != nil {
				return nil, err
			}
		}
	}

	return desiredRows, nil
}

// syncPagePartialColumns synchronizes the PagePartialRowColumns for a given PagePartialRow based on the provided DTO columns.
func syncPagePartialColumns(tx *gorm.DB, partialID, rowID uint, existingColumns []models.PagePartialRowColumn, dtoColumns []requests.UpdatePagePartialRowColumn, depth int) ([]models.PagePartialRowColumn, error) {
	existingByID := make(map[uint]*models.PagePartialRowColumn, len(existingColumns))
	for i := range existingColumns {
		existingByID[existingColumns[i].ID] = &existingColumns[i]
	}

	desiredColumns := make([]models.PagePartialRowColumn, 0, len(dtoColumns))
	desiredIDs := make(map[uint]bool, len(dtoColumns))

	for i := range dtoColumns {
		dtoCol := dtoColumns[i]
		var col *models.PagePartialRowColumn
		if dtoCol.ID != nil {
			col = existingByID[*dtoCol.ID]
		}
		if col == nil {
			col = &models.PagePartialRowColumn{}
		}

		col.PagePartialRowID = rowID
		col.Position = utils.UintOrZero(dtoCol.Position)
		col.ModuleID = utils.NewNullUInt(dtoCol.ModuleID)
		col.Cols = dtoCol.Cols
		col.Xxl = utils.NewNullInt16(dtoCol.Xxl)
		col.Xl = utils.NewNullInt16(dtoCol.Xl)
		col.Lg = utils.NewNullInt16(dtoCol.Lg)
		col.Md = utils.NewNullInt16(dtoCol.Md)
		col.Sm = utils.NewNullInt16(dtoCol.Sm)
		col.Xs = utils.NewNullInt16(dtoCol.Xs)
		col.Offset = utils.NewNullInt16(dtoCol.Offset)
		col.OffsetXxl = utils.NewNullInt16(dtoCol.OffsetXxl)
		col.OffsetXl = utils.NewNullInt16(dtoCol.OffsetXl)
		col.OffsetLg = utils.NewNullInt16(dtoCol.OffsetLg)
		col.OffsetMd = utils.NewNullInt16(dtoCol.OffsetMd)
		col.OffsetSm = utils.NewNullInt16(dtoCol.OffsetSm)
		col.Order = utils.NewNullInt16(dtoCol.Order)
		col.OrderXxl = utils.NewNullInt16(dtoCol.OrderXxl)
		col.OrderXl = utils.NewNullInt16(dtoCol.OrderXl)
		col.OrderLg = utils.NewNullInt16(dtoCol.OrderLg)
		col.OrderMd = utils.NewNullInt16(dtoCol.OrderMd)
		col.OrderSm = utils.NewNullInt16(dtoCol.OrderSm)
		col.AlignSelf = utils.NewNullString(dtoCol.AlignSelf)
		col.Content = utils.NewNullString(dtoCol.Content)

		if col.ID != 0 {
			if err := tx.Model(col).
				Clauses(clause.Returning{Columns: []clause.Column{{Name: "updated_at"}}}).
				Updates(col).Error; err != nil {
				return nil, err
			}
		} else {
			if err := tx.Create(col).Error; err != nil {
				return nil, err
			}
		}

		existingNestedRows := make([]models.PagePartialRow, 0)
		if err := tx.Model(&models.PagePartialRow{}).
			Joins("JOIN page_partial_row_column_rows ON page_partial_row_column_rows.row_id = page_partial_rows.id").
			Where("page_partial_row_column_rows.column_id = ?", col.ID).
			Find(&existingNestedRows).Error; err != nil {
			return nil, err
		}

		nestedRows, err := syncPagePartialRows(tx, partialID, &col.ID, existingNestedRows, dtoCol.Rows, depth+1)
		if err != nil {
			return nil, err
		}

		col.PagePartialRows = nestedRows
		desiredColumns = append(desiredColumns, *col)
		desiredIDs[col.ID] = true
	}

	for i := range existingColumns {
		if _, ok := desiredIDs[existingColumns[i].ID]; !ok {
			if err := tx.Delete(&existingColumns[i]).Error; err != nil {
				return nil, err
			}
		}
	}

	return desiredColumns, nil
}

// DeletePage method to delete a page.
func DeletePage(menuItemID uint, locale string) error {
	err := database.Pg.Delete(&models.Page{MenuItemID: menuItemID, Locale: locale}).Error
	if err == nil {
		_ = deletePageFromCache(menuItemID, locale)
	}

	return err
}

// DeletePagePartial method to delete a page partial by its ID.
func DeletePagePartial(menuItemID uint, locale string, partialID uint) error {
	err := database.Pg.Delete(&models.PagePartial{}, partialID).Error
	if err == nil {
		_ = deletePageFromCache(menuItemID, locale)
	}

	return err
}

// RestorePage method to restore a deleted page.
func RestorePage(menuItemID uint, locale string) error {
	return database.Pg.Unscoped().
		Model(&models.Page{}).
		Where("menu_item_id = ? AND locale = ?", menuItemID, locale).
		Update("deleted_at", nil).Error
}

// RestorePagePartial method to restore a deleted page partial by its ID.
func RestorePagePartial(menuItemID uint, locale string, partialID uint) error {
	err := database.Pg.Unscoped().
		Model(&models.PagePartial{}).
		Where("id = ?", partialID).
		Update("deleted_at", nil).Error
	if err == nil {
		_ = deletePageFromCache(menuItemID, locale)
	}

	return err
}

// shared row ordering for root partial rows; exclude rows linked as nested column rows.
func preloadRootPagePartialRows(db *gorm.DB) *gorm.DB {
	return db.
		Where("NOT EXISTS (SELECT 1 FROM page_partial_row_column_rows prcr WHERE prcr.row_id = page_partial_rows.id)").
		Order("position asc")
}

// shared row ordering for nested rows loaded through column relations.
func preloadNestedPagePartialRows(db *gorm.DB) *gorm.DB {
	return db.Order("position asc")
}

// shared column ordering + module preload for partial tree preloads.
func preloadPagePartialColumns(db *gorm.DB) *gorm.DB {
	return db.Preload("Module").Order("position asc")
}

// preloadPagePartialTree loads rows/columns recursively up to MaxPagePartialTreeDepth.
func preloadPagePartialTree(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Rows", preloadRootPagePartialRows).
		Preload("Rows.Columns", preloadPagePartialColumns).
		Preload("Rows.Columns.PagePartialRows", preloadNestedPagePartialRows).
		Preload("Rows.Columns.PagePartialRows.Columns", preloadPagePartialColumns).
		Preload("Rows.Columns.PagePartialRows.Columns.PagePartialRows", preloadNestedPagePartialRows).
		Preload("Rows.Columns.PagePartialRows.Columns.PagePartialRows.Columns", preloadPagePartialColumns).
		Preload("Rows.Columns.PagePartialRows.Columns.PagePartialRows.Columns.PagePartialRows", preloadNestedPagePartialRows).
		Preload("Rows.Columns.PagePartialRows.Columns.PagePartialRows.Columns.PagePartialRows.Columns", preloadPagePartialColumns)
}

// ensureRowColumnRelation ensures that a relation exists between a row and a column in the page_partial_row_column_rows table.
func ensureRowColumnRelation(tx *gorm.DB, columnID, rowID uint) error {
	relation := models.PagePartialRowColumnRow{ColumnID: columnID, RowID: rowID}
	return tx.Where("column_id = ? AND row_id = ?", columnID, rowID).FirstOrCreate(&relation).Error
}

// getPageCacheKey gets the key for the cache.
func getPageCacheKey(menuItemID uint, locale string) string {
	return fmt.Sprintf("pages:%d:%s", menuItemID, locale)
}

// isPageInCache checks if the page exists in the cache.
func isPageInCache(menuItemID uint, locale string) (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(getPageCacheKey(menuItemID, locale)).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// getPageFromCache gets the page from the cache.
func getPageFromCache(menuItemID uint, locale string) (*models.Page, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(getPageCacheKey(menuItemID, locale)).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var page models.Page
	if err := json.Unmarshal([]byte(value), &page); err != nil {
		return nil, err
	}

	return &page, nil
}

// setPageToCache sets the page to the cache.
func setPageToCache(menuItemID uint, locale string, page *models.Page) error {
	value, err := json.Marshal(page)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(getPageCacheKey(menuItemID, locale)).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deletePageFromCache deletes existing page from the cache.
func deletePageFromCache(menuItemID uint, locale string) error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(getPageCacheKey(menuItemID, locale)).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deletePagesFromCacheByMenuItemID deletes all pages related to a menu item from the cache by the menu item ID.
func deletePagesFromCacheByMenuItemID(menuItemID uint) error {
	var locales []string
	if result := database.Pg.Model(&models.Page{}).Where("menu_item_id = ?", menuItemID).Pluck("locale", &locales); result.Error != nil {
		return result.Error
	}

	for i := range locales {
		if err := deletePageFromCache(menuItemID, locales[i]); err != nil {
			return err
		}
	}

	return nil
}
