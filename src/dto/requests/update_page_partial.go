package requests

import "time"

type UpdatePagePartial struct {
	Name      string                 `json:"name" validate:"required"`
	UpdatedAt time.Time              `json:"updatedAt" validate:"required"`
	Rows      []UpdatePagePartialRow `json:"rows" validate:"required,min=1,dive"`
}
