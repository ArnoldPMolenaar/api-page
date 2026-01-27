package services

import (
	"api-page/main/src/database"
	"api-page/main/src/dto/requests"
	"api-page/main/src/dto/responses"
	"api-page/main/src/enums"
	"api-page/main/src/models"
	"database/sql"
	"sort"

	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/gofiber/fiber/v2"
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

// GetMenuByID method to get a menu by ID.
func GetMenuByID(menuID uint) (*models.Menu, error) {
	menu := &models.Menu{}

	if result := database.Pg.
		Preload("MenuItemRelations", func(db *gorm.DB) *gorm.DB {
			return db.Preload("MenuItemParent", func(db *gorm.DB) *gorm.DB { return db.Preload("Indexing") }).
				Preload("MenuItemChild", func(db *gorm.DB) *gorm.DB { return db.Preload("Indexing") }).
				Order("menu_item_parent_id NULLS FIRST").
				Order("position ASC")
		}).
		Find(menu, "id = ?", menuID); result.Error != nil {
		return nil, result.Error
	}

	return menu, nil
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

	return oldMenu, nil
}

// DeleteMenu method to delete a menu.
func DeleteMenu(menuID uint) error {
	err := database.Pg.Delete(&models.Menu{}, menuID).Error

	return err
}

// RestoreMenu method to restore a deleted menu.
func RestoreMenu(menuID uint) error {
	err := database.Pg.Unscoped().Model(&models.Menu{}).Where("id = ?", menuID).Update("deleted_at", nil).Error

	return err
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
