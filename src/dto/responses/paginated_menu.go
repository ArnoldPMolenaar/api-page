package responses

import (
	"api-page/main/src/models"
	"time"
)

type PaginatedMenu struct {
	ID        uint      `json:"id"`
	VersionID uint      `json:"versionId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SetPaginatedMenu method to set menu data from models.Menu{}.
func (m *PaginatedMenu) SetPaginatedMenu(menu *models.Menu) {
	m.ID = menu.ID
	m.VersionID = menu.VersionID
	m.Name = menu.Name
	m.CreatedAt = menu.CreatedAt
	m.UpdatedAt = menu.UpdatedAt
}
