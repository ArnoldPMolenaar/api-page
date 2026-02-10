package responses

import (
	"api-page/main/src/models"
)

type PublishedPagePartial struct {
	ID   uint                      `json:"id"`
	Name string                    `json:"name"`
	Rows []PublishedPagePartialRow `json:"rows"`
}

// SetPagePartial sets the PagePartial response from the models.PagePartial model.
func (ppp *PublishedPagePartial) SetPagePartial(partial *models.PagePartial) {
	ppp.ID = partial.ID
	ppp.Name = partial.Name

	ppp.Rows = make([]PublishedPagePartialRow, len(partial.Rows))
	for i := range partial.Rows {
		ppr := PublishedPagePartialRow{}
		ppr.SetPagePartialRow(&partial.Rows[i])
		ppp.Rows[i] = ppr
	}
}
