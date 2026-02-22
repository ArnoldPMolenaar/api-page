package requests

import "time"

type UpdateMenuItem struct {
	ID *uint `json:"id"`
	// Use a pointer to uint for Position to allow zero value and required validation.
	Position  *uint              `json:"position" validate:"required"`
	Name      string             `json:"name" validate:"required"`
	Icon      *string            `json:"icon"`
	UpdatedAt *time.Time         `json:"updatedAt"`
	EnabledAt *time.Time         `json:"enabledAt"`
	Indexing  []MenuItemIndexing `json:"indexing" validate:"required,min=1,dive"`
	Items     []UpdateMenuItem   `json:"items" validate:"dive"`
}
