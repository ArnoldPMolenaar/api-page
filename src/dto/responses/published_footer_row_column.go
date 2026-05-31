package responses

import (
	"api-page/main/src/models"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type PublishedFooterRowColumn struct {
	ID        uint                 `json:"id" `
	Position  uint                 `json:"position"`
	Cols      string               `json:"cols"`
	Xxl       *int16               `json:"xxl"`
	Xl        *int16               `json:"xl"`
	Lg        *int16               `json:"lg"`
	Md        *int16               `json:"md"`
	Sm        *int16               `json:"sm"`
	Xs        *int16               `json:"xs"`
	Offset    *int16               `json:"offset"`
	OffsetXxl *int16               `json:"offsetXxl"`
	OffsetXl  *int16               `json:"offsetXl"`
	OffsetLg  *int16               `json:"offsetLg"`
	OffsetMd  *int16               `json:"offsetMd"`
	OffsetSm  *int16               `json:"offsetSm"`
	Order     *int16               `json:"order"`
	OrderXxl  *int16               `json:"orderXxl"`
	OrderXl   *int16               `json:"orderXl"`
	OrderLg   *int16               `json:"orderLg"`
	OrderMd   *int16               `json:"orderMd"`
	OrderSm   *int16               `json:"orderSm"`
	AlignSelf *string              `json:"alignSelf"`
	Content   *string              `json:"content"`
	Module    *PublishedModule     `json:"module"`
	Rows      []PublishedFooterRow `json:"rows"`
}

// SetFooterRowColumn sets the FooterRowColumn response from the models.FooterRowColumn model.
func (pfrc *PublishedFooterRowColumn) SetFooterRowColumn(column *models.FooterRowColumn) {
	pfrc.ID = column.ID
	pfrc.Position = column.Position
	pfrc.Cols = column.Cols
	pfrc.Xxl = utils.PtrFromNullInt16(column.Xxl)
	pfrc.Xl = utils.PtrFromNullInt16(column.Xl)
	pfrc.Lg = utils.PtrFromNullInt16(column.Lg)
	pfrc.Md = utils.PtrFromNullInt16(column.Md)
	pfrc.Sm = utils.PtrFromNullInt16(column.Sm)
	pfrc.Xs = utils.PtrFromNullInt16(column.Xs)
	pfrc.Offset = utils.PtrFromNullInt16(column.Offset)
	pfrc.OffsetXxl = utils.PtrFromNullInt16(column.OffsetXxl)
	pfrc.OffsetXl = utils.PtrFromNullInt16(column.OffsetXl)
	pfrc.OffsetLg = utils.PtrFromNullInt16(column.OffsetLg)
	pfrc.OffsetMd = utils.PtrFromNullInt16(column.OffsetMd)
	pfrc.OffsetSm = utils.PtrFromNullInt16(column.OffsetSm)
	pfrc.Order = utils.PtrFromNullInt16(column.Order)
	pfrc.OrderXxl = utils.PtrFromNullInt16(column.OrderXxl)
	pfrc.OrderXl = utils.PtrFromNullInt16(column.OrderXl)
	pfrc.OrderLg = utils.PtrFromNullInt16(column.OrderLg)
	pfrc.OrderMd = utils.PtrFromNullInt16(column.OrderMd)
	pfrc.OrderSm = utils.PtrFromNullInt16(column.OrderSm)
	pfrc.AlignSelf = utils.PtrFromNullString(column.AlignSelf)
	pfrc.Content = utils.PtrFromNullString(column.Content)
	if column.Module != nil {
		var module PublishedModule
		module.SetModule(column.Module)
		pfrc.Module = &module
	}

	pfrc.Rows = make([]PublishedFooterRow, len(column.FooterRows))
	for i := range column.FooterRows {
		row := PublishedFooterRow{}
		row.SetFooterRow(&column.FooterRows[i])
		pfrc.Rows[i] = row
	}
}
