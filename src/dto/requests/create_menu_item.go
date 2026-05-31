package requests

import (
	"api-page/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type CreateMenuItem struct {
	// Use a pointer to uint for Position to allow zero value and required validation.
	Position  *uint              `json:"position" validate:"required"`
	Name      string             `json:"name" validate:"required"`
	Icon      *string            `json:"icon"`
	EnabledAt *time.Time         `json:"enabledAt"`
	Indexing  []MenuItemIndexing `json:"indexing" validate:"required,min=1,dive"`
	Items     []CreateMenuItem   `json:"items" validate:"dive"`
}

func (c *CreateMenuItem) SetMenuItemRelation(relation *models.MenuItemRelation) {
	indexing := make([]MenuItemIndexing, 0, len(relation.MenuItemChild.Indexing))
	for i := range relation.MenuItemChild.Indexing {
		menuItemIndexing := MenuItemIndexing{}
		menuItemIndexing.SetMenuItemIndexing(&relation.MenuItemChild.Indexing[i])
		indexing = append(indexing, menuItemIndexing)
	}

	c.Position = &relation.Position
	c.Name = relation.MenuItemChild.Name
	c.Icon = utils.PtrFromNullString(relation.MenuItemChild.Icon)
	c.EnabledAt = utils.PtrFromNullTime(relation.MenuItemChild.EnabledAt)
	c.Indexing = indexing
	c.Items = make([]CreateMenuItem, 0)
}
