package requests

type UpdatePagePartial struct {
	Name string                 `json:"name" validate:"required"`
	Rows []UpdatePagePartialRow `json:"rows" validate:"required,min=1,dive"`
}
