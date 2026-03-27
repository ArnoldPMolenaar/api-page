package requests

type UpdateFooter struct {
	Rows []UpdateFooterRow `json:"rows" validate:"required,dive"`
}
