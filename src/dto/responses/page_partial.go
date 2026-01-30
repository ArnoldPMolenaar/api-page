package responses

import (
	"api-page/main/src/models"
	"time"
)

type PagePartial struct {
	ID         uint             `json:"id"`
	MenuItemID uint             `json:"menuItemId"`
	Locale     string           `json:"locale"`
	Name       string           `json:"name"`
	CreatedAt  time.Time        `json:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt"`
	Rows       []PagePartialRow `json:"rows"`
}

// SetPagePartial sets the PagePartial response from the models.PagePartial model.
func (pp *PagePartial) SetPagePartial(partial *models.PagePartial) {
	pp.ID = partial.ID
	pp.MenuItemID = partial.MenuItemID
	pp.Locale = partial.Locale
	pp.Name = partial.Name
	pp.CreatedAt = partial.CreatedAt
	pp.UpdatedAt = partial.UpdatedAt

	pp.Rows = make([]PagePartialRow, len(partial.Rows))
	for i := range partial.Rows {
		pr := PagePartialRow{}
		pr.SetPagePartialRow(&partial.Rows[i])
		pp.Rows[i] = pr
	}
}
