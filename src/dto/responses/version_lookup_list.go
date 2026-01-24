package responses

import "api-page/main/src/models"

type VersionLookupList struct {
	Versions []VersionLookup `json:"versions"`
}

// SetVersionLookupList sets the list of version lookups.
func (vll *VersionLookupList) SetVersionLookupList(versions *[]models.Version) {
	vll.Versions = make([]VersionLookup, len(*versions))
	for i, version := range *versions {
		var vl VersionLookup
		vl.SetVersionLookup(&version)
		vll.Versions[i] = vl
	}
}
