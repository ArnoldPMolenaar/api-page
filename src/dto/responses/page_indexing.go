package responses

import "api-page/main/src/models"

type PageIndexing struct {
	Option string  `json:"option"`
	Value  *string `json:"value"`
}

// SetPageIndexing sets the PageIndexing response from the models.PageIndexing model.
func (pi *PageIndexing) SetPageIndexing(indexing *models.PageIndexing) {
	pi.Option = indexing.Option.String()

	if indexing.Value.Valid {
		pi.Value = &indexing.Value.String
	}
}
