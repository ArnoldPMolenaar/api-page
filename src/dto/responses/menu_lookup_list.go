package responses

import "api-page/main/src/models"

type MenuLookupList struct {
	Menus []MenuLookup `json:"menus"`
}

// SetMenuLookupList sets the list of menu lookups.
func (mll *MenuLookupList) SetMenuLookupList(menus *[]models.Menu) {
	mll.Menus = make([]MenuLookup, len(*menus))
	for i := range *menus {
		var ml MenuLookup
		ml.SetMenuLookup(&(*menus)[i])
		mll.Menus[i] = ml
	}
}
