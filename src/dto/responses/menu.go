package responses

import (
	"api-page/main/src/models"
	"time"
)

type Menu struct {
	ID        uint       `json:"id"`
	VersionID uint       `json:"versionId"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	Items     []MenuItem `json:"items"`
}

// SetMenu sets the Menu response from the models.Menu model.
func (m *Menu) SetMenu(menu *models.Menu) {
	m.ID = menu.ID
	m.VersionID = menu.VersionID
	m.Name = menu.Name
	m.CreatedAt = menu.CreatedAt
	m.UpdatedAt = menu.UpdatedAt

	m.Items = make([]MenuItem, 0)
	for i, _ := range menu.MenuItemRelations {
		relation := &menu.MenuItemRelations[i]

		if !relation.MenuItemParentID.Valid {
			menuItem := MenuItem{}
			menuItem.SetMenuItem(&relation.MenuItemChild, relation.Position)
			m.Items = append(m.Items, menuItem)
			continue
		}

		parentItem := m.findParentItem(relation.MenuItemParentID.V, &m.Items)
		if parentItem != nil {
			menuItem := MenuItem{}
			menuItem.SetMenuItem(&relation.MenuItemChild, relation.Position)
			parentItem.Items = append(parentItem.Items, menuItem)
		}
	}
}

// findParentItem recursively finds a parent MenuItem by ID.
func (m *Menu) findParentItem(parentID uint, items *[]MenuItem) *MenuItem {
	for i, _ := range *items {
		item := (*items)[i]
		if item.ID == parentID {
			return &(*items)[i]
		}
		if len(item.Items) > 0 {
			if foundItem := m.findParentItem(parentID, &item.Items); foundItem != nil {
				return foundItem
			}
		}
	}

	return nil
}
