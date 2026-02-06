package requests

import "encoding/json"

// CreateModule represents the request payload for creating a new module.
type CreateModule struct {
	AppName  string          `json:"appName" validate:"required"`
	Type     string          `json:"type" validate:"required"`
	Name     string          `json:"name" validate:"required"`
	Settings json.RawMessage `json:"settings" validate:"required,validjson"`
}
