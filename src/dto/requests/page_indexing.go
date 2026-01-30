package requests

type PageIndexing struct {
	Option string  `json:"option" validate:"required"`
	Value  *string `json:"value"`
}
