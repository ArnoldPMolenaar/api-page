package responses

import (
	"api-page/main/src/models"
	"time"
)

type PaginatedModule struct {
	ID        uint      `json:"id"`
	AppName   string    `json:"appName"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SetPaginatedModule method to set module data from models.Module{}.
func (m *PaginatedModule) SetPaginatedModule(module *models.Module) {
	m.ID = module.ID
	m.AppName = module.AppName
	m.Type = module.Type
	m.Name = module.Name
	m.CreatedAt = module.CreatedAt
	m.UpdatedAt = module.UpdatedAt
}
