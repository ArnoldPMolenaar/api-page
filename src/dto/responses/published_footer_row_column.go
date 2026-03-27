package responses

import (
	"api-page/main/src/models"
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
	if column.Xxl.Valid {
		pfrc.Xxl = &column.Xxl.Int16
	}
	if column.Xl.Valid {
		pfrc.Xl = &column.Xl.Int16
	}
	if column.Lg.Valid {
		pfrc.Lg = &column.Lg.Int16
	}
	if column.Md.Valid {
		pfrc.Md = &column.Md.Int16
	}
	if column.Sm.Valid {
		pfrc.Sm = &column.Sm.Int16
	}
	if column.Xs.Valid {
		pfrc.Xs = &column.Xs.Int16
	}
	if column.Offset.Valid {
		pfrc.Offset = &column.Offset.Int16
	}
	if column.OffsetXxl.Valid {
		pfrc.OffsetXxl = &column.OffsetXxl.Int16
	}
	if column.OffsetXl.Valid {
		pfrc.OffsetXl = &column.OffsetXl.Int16
	}
	if column.OffsetLg.Valid {
		pfrc.OffsetLg = &column.OffsetLg.Int16
	}
	if column.OffsetMd.Valid {
		pfrc.OffsetMd = &column.OffsetMd.Int16
	}
	if column.OffsetSm.Valid {
		pfrc.OffsetSm = &column.OffsetSm.Int16
	}
	if column.Order.Valid {
		pfrc.Order = &column.Order.Int16
	}
	if column.OrderXxl.Valid {
		pfrc.OrderXxl = &column.OrderXxl.Int16
	}
	if column.OrderXl.Valid {
		pfrc.OrderXl = &column.OrderXl.Int16
	}
	if column.OrderLg.Valid {
		pfrc.OrderLg = &column.OrderLg.Int16
	}
	if column.OrderMd.Valid {
		pfrc.OrderMd = &column.OrderMd.Int16
	}
	if column.OrderSm.Valid {
		pfrc.OrderSm = &column.OrderSm.Int16
	}
	if column.AlignSelf.Valid {
		pfrc.AlignSelf = &column.AlignSelf.String
	}
	if column.Content.Valid {
		pfrc.Content = &column.Content.String
	}
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
