package responses

import "api-page/main/src/models"

type Footer struct {
	Rows []FooterRow `json:"rows"`
}

// SetFooter to bind the rows from models.FooterRow.
func (f *Footer) SetFooter(rows *[]models.FooterRow) {
	f.Rows = make([]FooterRow, len(*rows))
	for i, row := range *rows {
		f.Rows[i] = FooterRow{}
		f.Rows[i].SetFooterRow(&row)
	}
}
