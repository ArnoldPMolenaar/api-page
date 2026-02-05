package responses

import "api-page/main/src/models"

type PublishedMenuList struct {
	Menus []PublishedMenu `json:"menus"`
}

// SetMenuList sets the list of menus.
func (pml *PublishedMenuList) SetMenuList(menus *[]models.Menu) {
	pml.Menus = make([]PublishedMenu, len(*menus))
	for i, menu := range *menus {
		var pm PublishedMenu
		pm.SetMenu(&menu)
		pml.Menus[i] = pm
	}
}
