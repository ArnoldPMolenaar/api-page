package responses

import (
	"api-page/main/src/models"
	"time"
)

type MenuItem struct {
	ID        uint               `json:"id"`
	VersionID uint               `json:"versionId"`
	Position  uint               `json:"position"`
	Name      string             `json:"name"`
	Icon      *string            `json:"icon"`
	EnabledAt *time.Time         `json:"enabledAt"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
	Indexing  []MenuItemIndexing `json:"indexing"`
	Items     []MenuItem         `json:"items"`
}

// SetMenuItem sets the MenuItem response from the models.MenuItem model.
func (mi *MenuItem) SetMenuItem(menuItem *models.MenuItem, position uint) {
	mi.ID = menuItem.ID
	mi.VersionID = menuItem.VersionID
	mi.Position = position
	mi.Name = menuItem.Name

	if menuItem.Icon.Valid {
		mi.Icon = &menuItem.Icon.String
	}

	if menuItem.EnabledAt.Valid {
		mi.EnabledAt = &menuItem.EnabledAt.Time
	}

	mi.CreatedAt = menuItem.CreatedAt
	mi.UpdatedAt = menuItem.UpdatedAt

	mi.Indexing = make([]MenuItemIndexing, len(menuItem.Indexing))
	for i, indexing := range menuItem.Indexing {
		mi.Indexing[i] = MenuItemIndexing{}
		mi.Indexing[i].SetMenuIndexing(&indexing)
	}

	mi.Items = make([]MenuItem, 0)
}
