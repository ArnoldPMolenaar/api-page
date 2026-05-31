package requests

import (
	"api-page/main/src/models"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type MenuItemIndexing struct {
	Option string  `json:"option" validate:"required"`
	Value  *string `json:"value"`
}

func (m *MenuItemIndexing) SetMenuItemIndexing(indexing *models.MenuItemIndexing) {
	m.Option = indexing.Option.String()
	m.Value = utils.PtrFromNullString(indexing.Value)
}
