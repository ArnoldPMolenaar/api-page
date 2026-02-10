package responses

import (
	"api-page/main/src/models"
	"encoding/json"
)

type PublishedModule struct {
	Type     string          `json:"type"`
	Settings json.RawMessage `json:"settings"`
}

// SetModule sets the Module response from the models.Module model.
func (pm *PublishedModule) SetModule(module *models.Module) {
	pm.Type = module.Type
	pm.Settings = json.RawMessage(module.Settings)
}
