package requests

import (
	"encoding/json"
	"time"
)

// UpdateModule represents the request payload for updating an existing module.
type UpdateModule struct {
	Type      string          `json:"type" validate:"required"`
	Name      string          `json:"name" validate:"required"`
	Settings  json.RawMessage `json:"settings" validate:"required,validjson"`
	UpdatedAt time.Time       `json:"updatedAt" validate:"required"`
}
