package requests

import "time"

// UpdateMenu represents the request payload for updating a menu.
type UpdateMenu struct {
	Name      string           `json:"name" validate:"required"`
	UpdatedAt time.Time        `json:"updatedAt" validate:"required"`
	Items     []UpdateMenuItem `json:"items" validate:"required,min=1,dive"`
}
