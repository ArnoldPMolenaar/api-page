package requests

import (
	"api-page/main/src/models"
	"time"
)

type UpdatePagePartial struct {
	Name      string                 `json:"name" validate:"required"`
	UpdatedAt time.Time              `json:"updatedAt" validate:"required"`
	Rows      []UpdatePagePartialRow `json:"rows" validate:"required,min=1,dive"`
}

func (u *UpdatePagePartial) SetPagePartial(partial models.PagePartial, partialID uint) {
	rows := make([]UpdatePagePartialRow, 0, len(partial.Rows))
	for i := range partial.Rows {
		row := UpdatePagePartialRow{}
		row.SetPagePartialRow(partial.Rows[i], partialID)
		rows = append(rows, row)
	}

	u.Name = partial.Name
	u.UpdatedAt = partial.UpdatedAt
	u.Rows = rows
}
