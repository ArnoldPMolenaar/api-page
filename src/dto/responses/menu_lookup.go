package responses

import "api-page/main/src/models"

type MenuLookup struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// SetMenuLookup sets the menu lookup fields from a Menu model.
func (ml *MenuLookup) SetMenuLookup(menu *models.Menu) {
	ml.ID = menu.ID
	ml.Name = menu.Name
}
