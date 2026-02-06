package responses

import "api-page/main/src/models"

type ModuleLookupList struct {
	Modules []ModuleLookup `json:"modules"`
}

// SetModuleLookupList sets the list of module lookups.
func (mll *ModuleLookupList) SetModuleLookupList(modules *[]models.Module) {
	mll.Modules = make([]ModuleLookup, len(*modules))
	for i, module := range *modules {
		var ml ModuleLookup
		ml.SetModuleLookup(&module)
		mll.Modules[i] = ml
	}
}
