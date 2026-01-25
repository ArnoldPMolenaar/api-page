package requests

type MenuItemIndexing struct {
	Option string  `json:"option" validate:"required"`
	Value  *string `json:"value"`
}
