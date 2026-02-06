package responses

import (
	"api-page/main/src/models"
	"encoding/json"
	"time"
)

type Module struct {
	ID        uint            `json:"id"`
	AppName   string          `json:"appName"`
	Type      string          `json:"type"`
	Name      string          `json:"name"`
	Settings  json.RawMessage `json:"settings"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

// SetModule sets the Module response from the models.Module model.
func (m *Module) SetModule(module *models.Module) {
	m.ID = module.ID
	m.AppName = module.AppName
	m.Type = module.Type
	m.Name = module.Name
	m.Settings = json.RawMessage(module.Settings)
	m.CreatedAt = module.CreatedAt
	m.UpdatedAt = module.UpdatedAt
}
