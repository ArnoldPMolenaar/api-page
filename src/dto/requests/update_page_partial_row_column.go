package requests

import "time"

type UpdatePagePartialRowColumn struct {
	ID       *uint `json:"id" `
	RowID    *uint `json:"rowId"`
	ModuleID *uint `json:"moduleId"`
	// Use a pointer to uint for Position to allow zero value and required validation.
	Position  *uint      `json:"position" validate:"required"`
	Cols      string     `json:"cols" validate:"required"`
	Xxl       *int16     `json:"xxl"`
	Xl        *int16     `json:"xl"`
	Lg        *int16     `json:"lg"`
	Md        *int16     `json:"md"`
	Sm        *int16     `json:"sm"`
	Xs        *int16     `json:"xs"`
	Offset    *int16     `json:"offset"`
	OffsetXxl *int16     `json:"offsetXxl"`
	OffsetXl  *int16     `json:"offsetXl"`
	OffsetLg  *int16     `json:"offsetLg"`
	OffsetMd  *int16     `json:"offsetMd"`
	OffsetSm  *int16     `json:"offsetSm"`
	Order     *int16     `json:"order"`
	OrderXxl  *int16     `json:"orderXxl"`
	OrderXl   *int16     `json:"orderXl"`
	OrderLg   *int16     `json:"orderLg"`
	OrderMd   *int16     `json:"orderMd"`
	OrderSm   *int16     `json:"orderSm"`
	AlignSelf *string    `json:"alignSelf"`
	Content   *string    `json:"content"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
