package responses

import "api-page/main/src/models"

type VersionLookup struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// SetVersionLookup sets the version lookup fields from a Version model.
func (vl *VersionLookup) SetVersionLookup(version *models.Version) {
	vl.ID = version.ID
	vl.Name = version.Name
}
