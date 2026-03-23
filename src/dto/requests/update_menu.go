package requests

import "time"

// UpdateMenu represents the request payload for updating a menu.
type UpdateMenu struct {
	Name      string           `json:"name" validate:"required"`
	Depth     *uint8           `json:"depth"`
	UpdatedAt time.Time        `json:"updatedAt" validate:"required"`
	Items     []UpdateMenuItem `json:"items" validate:"required,min=1,dive"`
}
