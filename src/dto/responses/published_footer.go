package responses

import "api-page/main/src/models"

type PublishedFooter struct {
	Rows []PublishedFooterRow `json:"rows"`
}

// SetFooter to bind the rows from models.FooterRow.
func (pfr *PublishedFooter) SetFooter(rows *[]models.FooterRow) {
	pfr.Rows = make([]PublishedFooterRow, len(*rows))
	for i, row := range *rows {
		pfr.Rows[i] = PublishedFooterRow{}
		pfr.Rows[i].SetFooterRow(&row)
	}
}
