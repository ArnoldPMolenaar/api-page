package requests

import "time"

type CreateMenuItem struct {
	Position  uint               `json:"position" validate:"required"`
	Name      string             `json:"name" validate:"required"`
	Icon      *string            `json:"icon"`
	EnabledAt *time.Time         `json:"enabledAt"`
	Indexing  []MenuItemIndexing `json:"indexing" validate:"required,min=1,dive"`
	Items     []CreateMenuItem   `json:"items"`
}
