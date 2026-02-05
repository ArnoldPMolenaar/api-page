package responses

import (
	"api-page/main/src/models"
)

type PublishedMenu struct {
	ID    uint                `json:"id"`
	Name  string              `json:"name"`
	Items []PublishedMenuItem `json:"items"`
}

// SetMenu sets the Menu response from the models.Menu model.
func (pm *PublishedMenu) SetMenu(menu *models.Menu) {
	pm.ID = menu.ID
	pm.Name = menu.Name

	pm.Items = make([]PublishedMenuItem, 0)
	for i, _ := range menu.MenuItemRelations {
		relation := &menu.MenuItemRelations[i]

		if !relation.MenuItemParentID.Valid {
			menuItem := PublishedMenuItem{}
			menuItem.SetMenuItem(&relation.MenuItemChild, relation.Position)
			pm.Items = append(pm.Items, menuItem)
			continue
		}

		parentItem := pm.findParentItem(relation.MenuItemParentID.V, &pm.Items)
		if parentItem != nil {
			menuItem := PublishedMenuItem{}
			menuItem.SetMenuItem(&relation.MenuItemChild, relation.Position)
			parentItem.Items = append(parentItem.Items, menuItem)
		}
	}
}

// findParentItem recursively finds a parent MenuItem by ID.
func (pm *PublishedMenu) findParentItem(parentID uint, items *[]PublishedMenuItem) *PublishedMenuItem {
	for i, _ := range *items {
		item := (*items)[i]
		if item.ID == parentID {
			return &(*items)[i]
		}
		if len(item.Items) > 0 {
			if foundItem := pm.findParentItem(parentID, &item.Items); foundItem != nil {
				return foundItem
			}
		}
	}

	return nil
}
