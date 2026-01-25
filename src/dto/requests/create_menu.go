package requests

// CreateMenu represents the request payload for creating a new menu.
type CreateMenu struct {
	VersionID uint             `json:"versionId" validate:"required"`
	Name      string           `json:"name" validate:"required"`
	Items     []CreateMenuItem `json:"items" validate:"required,min=1,dive"`
}
