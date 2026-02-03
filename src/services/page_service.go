package services

import (
	"api-page/main/src/database"
	"api-page/main/src/dto/requests"
	"api-page/main/src/enums"
	"api-page/main/src/models"
	"api-page/main/src/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
		Preload("Partials", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Rows", func(db2 *gorm.DB) *gorm.DB {
				return db2.Preload("Columns", func(db3 *gorm.DB) *gorm.DB {
					return db3.Preload("Module").Order("position asc")
				}).Order("position asc")
			})
		}).
		FirstOrCreate(page, page)
	if result.Error != nil {
		return nil, result.Error
	}

	return page, nil
}

// GetPartialByID retrieves a PagePartial by its ID, including its associated rows and columns.
func GetPartialByID(partialID uint) (*models.PagePartial, error) {
	partial := &models.PagePartial{}

	if result := database.Pg.
		Preload("Rows", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Columns", func(db2 *gorm.DB) *gorm.DB {
				return db2.Preload("Module").Order("position asc")
			}).Order("position asc")
		}).
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

	return page, nil
}

// UpdatePagePartial updates the given PagePartial and its associated rows and columns
// based on the data provided in the UpdatePagePartial request.
func UpdatePagePartial(partial *models.PagePartial, dtoPartial *requests.UpdatePagePartial) (*models.PagePartial, error) {
	partial.Name = dtoPartial.Name

	if err := database.Pg.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&partial).
			Clauses(clause.Returning{Columns: []clause.Column{{Name: "updated_at"}}}).
			Updates(partial).Error; err != nil {
			tx.Rollback()
			return err
		}

		existingRows := make([]models.PagePartialRow, len(partial.Rows))
		copy(existingRows, partial.Rows)
		partial.Rows = make([]models.PagePartialRow, 0)

		desiredRowIDs := make(map[uint]bool)
		for i := range dtoPartial.Rows {
			dtoRow := dtoPartial.Rows[i]
			var row *models.PagePartialRow
			if dtoRow.ID != nil {
				for j := range existingRows {
					if existingRows[j].ID == *dtoRow.ID {
						row = &existingRows[j]
						break
					}
				}
			}
			if row == nil {
				row = &models.PagePartialRow{
					PagePartialID: partial.ID,
				}
			}

			row.Position = *dtoRow.Position
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
			row.Columns = make([]models.PagePartialRowColumn, 0)

			if row.ID != 0 {
				if err := tx.Model(&row).
					Clauses(clause.Returning{Columns: []clause.Column{{Name: "updated_at"}}}).
					Updates(row).Error; err != nil {
					tx.Rollback()
					return err
				}
			} else {
				if err := tx.FirstOrCreate(row, row).Error; err != nil {
					tx.Rollback()
					return err
				}
			}

			desiredRowIDs[row.ID] = true

			existingColumns := make([]models.PagePartialRowColumn, 0)
			if err := tx.Where("page_partial_row_id = ?", row.ID).Find(&existingColumns).Error; err != nil {
				tx.Rollback()
				return err
			}

			desiredColumnIDs := make(map[uint]bool)
			for i := range dtoRow.Columns {
				dtoCol := dtoRow.Columns[i]
				var col *models.PagePartialRowColumn
				if dtoCol.ID != nil {
					for j := range existingColumns {
						if existingColumns[j].ID == *dtoCol.ID {
							col = &existingColumns[j]
							break
						}
					}
				}
				if col == nil {
					col = &models.PagePartialRowColumn{
						PagePartialRowID: row.ID,
					}
				}

				col.Position = *dtoCol.Position
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
					if err := tx.Model(&col).
						Clauses(clause.Returning{Columns: []clause.Column{{Name: "updated_at"}}}).
						Updates(col).Error; err != nil {
						tx.Rollback()
						return err
					}
				} else {
					if err := tx.FirstOrCreate(col, col).Error; err != nil {
						tx.Rollback()
						return err
					}
				}

				desiredColumnIDs[col.ID] = true
				row.Columns = append(row.Columns, *col)
			}

			partial.Rows = append(partial.Rows, *row)

			for i := range existingColumns {
				if _, ok := desiredColumnIDs[existingColumns[i].ID]; !ok {
					if err := tx.Delete(&existingColumns[i]).Error; err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}

		for i := range existingRows {
			if _, ok := desiredRowIDs[existingRows[i].ID]; !ok {
				if err := tx.Delete(&existingRows[i]).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return partial, nil
}

// DeletePage method to delete a page.
func DeletePage(MenuItemID uint, Locale string) error {
	return database.Pg.Delete(&models.Page{MenuItemID: MenuItemID, Locale: Locale}).Error
}

// DeletePagePartial method to delete a page partial by its ID.
func DeletePagePartial(partialID uint) error {
	return database.Pg.Delete(&models.PagePartial{}, partialID).Error
}

// RestorePage method to restore a deleted page.
func RestorePage(MenuItemID uint, Locale string) error {
	return database.Pg.Unscoped().
		Model(&models.Page{}).
		Where("menu_item_id = ? AND locale = ?", MenuItemID, Locale).
		Update("deleted_at", nil).Error
}

// RestorePagePartial method to restore a deleted page partial by its ID.
func RestorePagePartial(partialID uint) error {
	return database.Pg.Unscoped().
		Model(&models.PagePartial{}).
		Where("id = ?", partialID).
		Update("deleted_at", nil).Error
}
