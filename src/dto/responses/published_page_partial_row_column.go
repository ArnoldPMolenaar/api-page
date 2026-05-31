package responses

import (
	"api-page/main/src/models"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type PublishedPagePartialRowColumn struct {
	ID        uint                      `json:"id" `
	Position  uint                      `json:"position"`
	Cols      string                    `json:"cols"`
	Xxl       *int16                    `json:"xxl"`
	Xl        *int16                    `json:"xl"`
	Lg        *int16                    `json:"lg"`
	Md        *int16                    `json:"md"`
	Sm        *int16                    `json:"sm"`
	Xs        *int16                    `json:"xs"`
	Offset    *int16                    `json:"offset"`
	OffsetXxl *int16                    `json:"offsetXxl"`
	OffsetXl  *int16                    `json:"offsetXl"`
	OffsetLg  *int16                    `json:"offsetLg"`
	OffsetMd  *int16                    `json:"offsetMd"`
	OffsetSm  *int16                    `json:"offsetSm"`
	Order     *int16                    `json:"order"`
	OrderXxl  *int16                    `json:"orderXxl"`
	OrderXl   *int16                    `json:"orderXl"`
	OrderLg   *int16                    `json:"orderLg"`
	OrderMd   *int16                    `json:"orderMd"`
	OrderSm   *int16                    `json:"orderSm"`
	AlignSelf *string                   `json:"alignSelf"`
	Content   *string                   `json:"content"`
	Module    *PublishedModule          `json:"module"`
	Rows      []PublishedPagePartialRow `json:"rows"`
}

// SetPagePartialRowColumn sets the PagePartialRowColumn response from the models.PagePartialRowColumn model.
func (ppprc *PublishedPagePartialRowColumn) SetPagePartialRowColumn(column *models.PagePartialRowColumn) {
	ppprc.ID = column.ID
	ppprc.Position = column.Position
	ppprc.Cols = column.Cols
	ppprc.Xxl = utils.PtrFromNullInt16(column.Xxl)
	ppprc.Xl = utils.PtrFromNullInt16(column.Xl)
	ppprc.Lg = utils.PtrFromNullInt16(column.Lg)
	ppprc.Md = utils.PtrFromNullInt16(column.Md)
	ppprc.Sm = utils.PtrFromNullInt16(column.Sm)
	ppprc.Xs = utils.PtrFromNullInt16(column.Xs)
	ppprc.Offset = utils.PtrFromNullInt16(column.Offset)
	ppprc.OffsetXxl = utils.PtrFromNullInt16(column.OffsetXxl)
	ppprc.OffsetXl = utils.PtrFromNullInt16(column.OffsetXl)
	ppprc.OffsetLg = utils.PtrFromNullInt16(column.OffsetLg)
	ppprc.OffsetMd = utils.PtrFromNullInt16(column.OffsetMd)
	ppprc.OffsetSm = utils.PtrFromNullInt16(column.OffsetSm)
	ppprc.Order = utils.PtrFromNullInt16(column.Order)
	ppprc.OrderXxl = utils.PtrFromNullInt16(column.OrderXxl)
	ppprc.OrderXl = utils.PtrFromNullInt16(column.OrderXl)
	ppprc.OrderLg = utils.PtrFromNullInt16(column.OrderLg)
	ppprc.OrderMd = utils.PtrFromNullInt16(column.OrderMd)
	ppprc.OrderSm = utils.PtrFromNullInt16(column.OrderSm)
	ppprc.AlignSelf = utils.PtrFromNullString(column.AlignSelf)
	ppprc.Content = utils.PtrFromNullString(column.Content)
	if column.Module != nil {
		var module PublishedModule
		module.SetModule(column.Module)
		ppprc.Module = &module
	}

	ppprc.Rows = make([]PublishedPagePartialRow, len(column.PagePartialRows))
	for i := range column.PagePartialRows {
		row := PublishedPagePartialRow{}
		row.SetPagePartialRow(&column.PagePartialRows[i])
		ppprc.Rows[i] = row
	}
}
