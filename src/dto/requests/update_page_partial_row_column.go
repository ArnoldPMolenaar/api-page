package requests

import (
	"api-page/main/src/models"
	"api-page/main/src/utils"
	"time"
)

type UpdatePagePartialRowColumn struct {
	ID       *uint `json:"id" `
	RowID    *uint `json:"rowId"`
	ModuleID *uint `json:"moduleId"`
	// Use a pointer to uint for Position to allow zero value and required validation.
	Position  *uint                  `json:"position" validate:"required"`
	Cols      string                 `json:"cols" validate:"required"`
	Xxl       *int16                 `json:"xxl"`
	Xl        *int16                 `json:"xl"`
	Lg        *int16                 `json:"lg"`
	Md        *int16                 `json:"md"`
	Sm        *int16                 `json:"sm"`
	Xs        *int16                 `json:"xs"`
	Offset    *int16                 `json:"offset"`
	OffsetXxl *int16                 `json:"offsetXxl"`
	OffsetXl  *int16                 `json:"offsetXl"`
	OffsetLg  *int16                 `json:"offsetLg"`
	OffsetMd  *int16                 `json:"offsetMd"`
	OffsetSm  *int16                 `json:"offsetSm"`
	Order     *int16                 `json:"order"`
	OrderXxl  *int16                 `json:"orderXxl"`
	OrderXl   *int16                 `json:"orderXl"`
	OrderLg   *int16                 `json:"orderLg"`
	OrderMd   *int16                 `json:"orderMd"`
	OrderSm   *int16                 `json:"orderSm"`
	AlignSelf *string                `json:"alignSelf"`
	Content   *string                `json:"content"`
	UpdatedAt *time.Time             `json:"updatedAt"`
	Rows      []UpdatePagePartialRow `json:"rows" validate:"dive"`
}

func (u *UpdatePagePartialRowColumn) SetPagePartialRowColumn(column models.PagePartialRowColumn, partialID uint) {
	rows := make([]UpdatePagePartialRow, 0, len(column.PagePartialRows))
	for i := range column.PagePartialRows {
		row := UpdatePagePartialRow{}
		row.SetPagePartialRow(column.PagePartialRows[i], partialID)
		rows = append(rows, row)
	}

	u.ModuleID = utils.PtrFromNullUint(column.ModuleID)
	u.Position = utils.PtrFromUint(column.Position)
	u.Cols = column.Cols
	u.Xxl = utils.PtrFromNullInt16(column.Xxl)
	u.Xl = utils.PtrFromNullInt16(column.Xl)
	u.Lg = utils.PtrFromNullInt16(column.Lg)
	u.Md = utils.PtrFromNullInt16(column.Md)
	u.Sm = utils.PtrFromNullInt16(column.Sm)
	u.Xs = utils.PtrFromNullInt16(column.Xs)
	u.Offset = utils.PtrFromNullInt16(column.Offset)
	u.OffsetXxl = utils.PtrFromNullInt16(column.OffsetXxl)
	u.OffsetXl = utils.PtrFromNullInt16(column.OffsetXl)
	u.OffsetLg = utils.PtrFromNullInt16(column.OffsetLg)
	u.OffsetMd = utils.PtrFromNullInt16(column.OffsetMd)
	u.OffsetSm = utils.PtrFromNullInt16(column.OffsetSm)
	u.Order = utils.PtrFromNullInt16(column.Order)
	u.OrderXxl = utils.PtrFromNullInt16(column.OrderXxl)
	u.OrderXl = utils.PtrFromNullInt16(column.OrderXl)
	u.OrderLg = utils.PtrFromNullInt16(column.OrderLg)
	u.OrderMd = utils.PtrFromNullInt16(column.OrderMd)
	u.OrderSm = utils.PtrFromNullInt16(column.OrderSm)
	u.AlignSelf = utils.PtrFromNullString(column.AlignSelf)
	u.Content = utils.PtrFromNullString(column.Content)
	u.Rows = rows
}
