package requests

import "time"

// UpdateVersion represents the request payload for updating a version.
type UpdateVersion struct {
	Name      string     `json:"name" validate:"required"`
	EnabledAt *time.Time `json:"enabledAt"`
	UpdatedAt time.Time  `json:"updatedAt" validate:"required"`
}
