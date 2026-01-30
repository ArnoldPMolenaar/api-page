package requests

type CreatePagePartial struct {
	Name string `json:"name" validate:"required"`
}
