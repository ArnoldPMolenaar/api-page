package responses

import "api-page/main/src/models"

type MenuItemIndexing struct {
	Option string  `json:"option"`
	Value  *string `json:"value"`
}

// SetMenuIndexing sets the MenuItemIndexing response from the models.MenuItemIndexing model.
func (mii *MenuItemIndexing) SetMenuIndexing(indexing *models.MenuItemIndexing) {
	mii.Option = indexing.Option.String()

	if indexing.Value.Valid {
		mii.Value = &indexing.Value.String
	}
}
