package services

import (
	"api-page/main/src/cache"
	"api-page/main/src/database"
	"api-page/main/src/dto/requests"
	"api-page/main/src/models"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/valkey-io/valkey-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const MaxRowTreeDepth = 4

// GetFooterByVersionID retrieves a Footer by its versionID.
func GetFooterByVersionID(versionID uint, locale string) (*[]models.FooterRow, error) {
	rows := make([]models.FooterRow, 0)

	if inCache, err := isFooterInCache(versionID, locale); err != nil {
		return nil, err
	} else if inCache {
		if cacheRows, err := getFooterFromCache(versionID, locale); err != nil {
			return nil, err
		} else if cacheRows != nil {
			rows = *cacheRows
		}
	}

	if len(rows) == 0 {
		if result := preloadFooterTree(database.Pg).
			Where("NOT EXISTS (SELECT 1 FROM footer_row_column_rows frcr WHERE frcr.row_id = footer_rows.id)").
			Order("position asc").
			Find(&rows, "version_id = ? AND locale = ?", versionID, locale); result.Error != nil {
			return nil, result.Error
		}
		_ = setFooterToCache(versionID, locale, &rows)
	}

	return &rows, nil
}

// UpdateFooter updates the given Footer and its associated rows and columns
// based on the data provided in the UpdateFooter request.
func UpdateFooter(versionID uint, locale string, footerRows *[]models.FooterRow, dtoFooter *requests.UpdateFooter) (*[]models.FooterRow, error) {
	var result *[]models.FooterRow

	if err := database.Pg.Transaction(func(tx *gorm.DB) error {
		var txErr error
		result, txErr = UpdateFooterWithTx(tx, versionID, locale, footerRows, dtoFooter)
		return txErr
	}); err != nil {
		return nil, err
	}

	_ = deleteFooterFromCache(versionID, locale)

	return result, nil
}

// UpdateFooterWithTx updates footer rows/columns using the provided transaction.
// It performs no transaction lifecycle control and no cache side effects.
func UpdateFooterWithTx(tx *gorm.DB, versionID uint, locale string, footerRows *[]models.FooterRow, dtoFooter *requests.UpdateFooter) (*[]models.FooterRow, error) {
	if tx == nil {
		return nil, gorm.ErrInvalidDB
	}

	result := make([]models.FooterRow, 0)

	existingRows := make([]models.FooterRow, len(*footerRows))
	copy(existingRows, *footerRows)

	rows, err := syncFooterRows(tx, versionID, locale, nil, existingRows, dtoFooter.Rows, 1)
	if err != nil {
		return nil, err
	}

	result = append(result, rows...)

	return &result, nil
}

// syncFooterRows synchronizes the FooterRows for a given Footer based on the provided DTO rows.
func syncFooterRows(tx *gorm.DB, versionID uint, locale string, parentColumnID *uint, existingRows []models.FooterRow, dtoRows []requests.UpdateFooterRow, depth int) ([]models.FooterRow, error) {
	if depth > MaxRowTreeDepth {
		return nil, fmt.Errorf("footer row depth exceeded max depth of %d", MaxRowTreeDepth)
	}

	existingByID := make(map[uint]*models.FooterRow, len(existingRows))
	for i := range existingRows {
		existingByID[existingRows[i].ID] = &existingRows[i]
	}

	desiredRows := make([]models.FooterRow, 0, len(dtoRows))
	desiredIDs := make(map[uint]bool, len(dtoRows))

	for i := range dtoRows {
		dtoRow := dtoRows[i]
		var row *models.FooterRow
		if dtoRow.ID != nil {
			row = existingByID[*dtoRow.ID]
		}
		if row == nil {
			row = &models.FooterRow{}
		}

		row.VersionID = versionID
		row.Locale = locale
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
			if err := ensureFooterRowColumnRelation(tx, *parentColumnID, row.ID); err != nil {
				return nil, err
			}
		}

		existingColumns := make([]models.FooterRowColumn, 0)
		if dtoRow.ID != nil && len(row.Columns) == 0 {
			if err := tx.Where("footer_row_id = ?", row.ID).Find(&existingColumns).Error; err != nil {
				return nil, err
			}
		} else {
			existingColumns = row.Columns
		}

		columns, err := syncFooterColumns(tx, versionID, row.ID, locale, existingColumns, dtoRow.Columns, depth)
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

// syncFooterColumns synchronizes the FooterRowColumns for a given FooterRow based on the provided DTO columns.
func syncFooterColumns(tx *gorm.DB, versionID, rowID uint, locale string, existingColumns []models.FooterRowColumn, dtoColumns []requests.UpdateFooterRowColumn, depth int) ([]models.FooterRowColumn, error) {
	existingByID := make(map[uint]*models.FooterRowColumn, len(existingColumns))
	for i := range existingColumns {
		existingByID[existingColumns[i].ID] = &existingColumns[i]
	}

	desiredColumns := make([]models.FooterRowColumn, 0, len(dtoColumns))
	desiredIDs := make(map[uint]bool, len(dtoColumns))

	for i := range dtoColumns {
		dtoCol := dtoColumns[i]
		var col *models.FooterRowColumn
		if dtoCol.ID != nil {
			col = existingByID[*dtoCol.ID]
		}
		if col == nil {
			col = &models.FooterRowColumn{}
		}

		col.FooterRowID = rowID
		col.Position = utils.UintOrZero(dtoCol.Position)
		col.ModuleID = utils.NewNull[uint](dtoCol.ModuleID)
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

		existingNestedRows := make([]models.FooterRow, 0)
		if err := tx.Model(&models.FooterRow{}).
			Joins("JOIN footer_row_column_rows ON footer_row_column_rows.row_id = footer_rows.id").
			Where("footer_row_column_rows.column_id = ?", col.ID).
			Find(&existingNestedRows).Error; err != nil {
			return nil, err
		}

		nestedRows, err := syncFooterRows(tx, versionID, locale, &col.ID, existingNestedRows, dtoCol.Rows, depth+1)
		if err != nil {
			return nil, err
		}

		col.FooterRows = nestedRows
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

// shared row ordering for nested rows loaded through column relations.
func preloadNestedFooterRows(db *gorm.DB) *gorm.DB {
	return db.Order("position asc")
}

// shared column ordering + module preload for row tree preloads.
func preloadFooterColumns(db *gorm.DB) *gorm.DB {
	return db.Preload("Module").Order("position asc")
}

// preloadFooterTree loads rows/columns recursively up to MaxRowTreeDepth.
func preloadFooterTree(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Columns", preloadFooterColumns).
		Preload("Columns.FooterRows", preloadNestedFooterRows).
		Preload("Columns.FooterRows.Columns", preloadFooterColumns).
		Preload("Columns.FooterRows.Columns.FooterRows", preloadNestedFooterRows).
		Preload("Columns.FooterRows.Columns.FooterRows.Columns", preloadFooterColumns).
		Preload("Columns.FooterRows.Columns.FooterRows.Columns.FooterRows", preloadNestedFooterRows).
		Preload("Columns.FooterRows.Columns.FooterRows.Columns.FooterRows.Columns", preloadFooterColumns)
}

// ensureFooterRowColumnRelation ensures that a relation exists between a row and a column in the footer_row_column_rows table.
func ensureFooterRowColumnRelation(tx *gorm.DB, columnID, rowID uint) error {
	relation := models.FooterRowColumnRow{ColumnID: columnID, RowID: rowID}
	return tx.Where("column_id = ? AND row_id = ?", columnID, rowID).FirstOrCreate(&relation).Error
}

// getFooterCacheKey gets the key for the cache.
func getFooterCacheKey(versionID uint, locale string) string {
	return fmt.Sprintf("footers:%d:%s", versionID, locale)
}

// isFooterInCache checks if the footer exists in the cache.
func isFooterInCache(versionID uint, locale string) (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(getFooterCacheKey(versionID, locale)).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// getFooterFromCache gets the footer from the cache.
func getFooterFromCache(versionID uint, locale string) (*[]models.FooterRow, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(getFooterCacheKey(versionID, locale)).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var rows []models.FooterRow
	if err := json.Unmarshal([]byte(value), &rows); err != nil {
		return nil, err
	}

	return &rows, nil
}

// setFooterToCache sets the footer rows to the cache.
func setFooterToCache(versionID uint, locale string, rows *[]models.FooterRow) error {
	value, err := json.Marshal(rows)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(getFooterCacheKey(versionID, locale)).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// deleteFooterFromCache deletes existing footer from the cache.
func deleteFooterFromCache(versionID uint, locale string) error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(getFooterCacheKey(versionID, locale)).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}
