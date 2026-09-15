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
func (pm *PublishedMenu) SetMenu(menu *models.Menu, pages *[]models.Page) {
	pm.ID = menu.ID
	pm.Name = menu.Name

	pm.Items = make([]PublishedMenuItem, 0)
	for i := range menu.MenuItemRelations {
		relation := &menu.MenuItemRelations[i]

		if !relation.MenuItemParentID.Valid {
			menuItem := PublishedMenuItem{}
			routeMenuItemIDs := appendRouteMenuItemID(nil, &relation.MenuItemChild)
			menuItem.SetMenuItem(&relation.MenuItemChild, relation.Position, pages, routeMenuItemIDs)
			pm.Items = append(pm.Items, menuItem)
			continue
		}

		parentItem := pm.findParentItem(relation.MenuItemParentID.V, &pm.Items)
		if parentItem != nil {
			menuItem := PublishedMenuItem{}
			routeMenuItemIDs := appendRouteMenuItemID(parentItem.routeMenuItemIDs, &relation.MenuItemChild)
			menuItem.SetMenuItem(&relation.MenuItemChild, relation.Position, pages, routeMenuItemIDs)
			parentItem.Items = append(parentItem.Items, menuItem)
		}
	}
}

// appendRouteMenuItemID copies the parent route path and appends the current item unless it is the homepage.
func appendRouteMenuItemID(parentRouteMenuItemIDs []uint, menuItem *models.MenuItem) []uint {
	routeMenuItemIDs := make([]uint, 0, len(parentRouteMenuItemIDs)+1)
	routeMenuItemIDs = append(routeMenuItemIDs, parentRouteMenuItemIDs...)
	if menuItem != nil && !menuItem.IsHome {
		routeMenuItemIDs = append(routeMenuItemIDs, menuItem.ID)
	}

	return routeMenuItemIDs
}

// findParentItem recursively finds a parent MenuItem by ID.
func (pm *PublishedMenu) findParentItem(parentID uint, items *[]PublishedMenuItem) *PublishedMenuItem {
	for i := range *items {
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
