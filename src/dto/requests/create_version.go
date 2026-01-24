package requests

import "time"

// CreateVersion represents the request payload for creating a new version.
type CreateVersion struct {
	AppName   string     `json:"appName" validate:"required"`
	Name      string     `json:"name" validate:"required"`
	EnabledAt *time.Time `json:"enabledAt"`
}
