package requests

import (
	"api-page/main/src/models"
	"api-page/main/src/utils"
)

type MenuItemIndexing struct {
	Option string  `json:"option" validate:"required"`
	Value  *string `json:"value"`
}

func (m *MenuItemIndexing) SetMenuItemIndexing(indexing models.MenuItemIndexing) {
	m.Option = indexing.Option.String()
	m.Value = utils.PtrFromNullString(indexing.Value)
}
