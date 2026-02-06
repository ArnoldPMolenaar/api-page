package responses

import "api-page/main/src/models"

type ModuleLookup struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// SetModuleLookup sets the module lookup fields from a Module model.
func (ml *ModuleLookup) SetModuleLookup(module *models.Module) {
	ml.ID = module.ID
	ml.Name = module.Name
}
