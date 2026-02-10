package responses

import (
	"api-page/main/src/models"
)

type PublishedPagePartialRowColumn struct {
	ID        uint             `json:"id" `
	Position  uint             `json:"position"`
	Cols      string           `json:"cols"`
	Xxl       *int16           `json:"xxl"`
	Xl        *int16           `json:"xl"`
	Lg        *int16           `json:"lg"`
	Md        *int16           `json:"md"`
	Sm        *int16           `json:"sm"`
	Xs        *int16           `json:"xs"`
	Offset    *int16           `json:"offset"`
	OffsetXxl *int16           `json:"offsetXxl"`
	OffsetXl  *int16           `json:"offsetXl"`
	OffsetLg  *int16           `json:"offsetLg"`
	OffsetMd  *int16           `json:"offsetMd"`
	OffsetSm  *int16           `json:"offsetSm"`
	Order     *int16           `json:"order"`
	OrderXxl  *int16           `json:"orderXxl"`
	OrderXl   *int16           `json:"orderXl"`
	OrderLg   *int16           `json:"orderLg"`
	OrderMd   *int16           `json:"orderMd"`
	OrderSm   *int16           `json:"orderSm"`
	AlignSelf *string          `json:"alignSelf"`
	Content   *string          `json:"content"`
	Module    *PublishedModule `json:"module"`
}

// SetPagePartialRowColumn sets the PagePartialRowColumn response from the models.PagePartialRowColumn model.
func (ppprc *PublishedPagePartialRowColumn) SetPagePartialRowColumn(column *models.PagePartialRowColumn) {
	ppprc.ID = column.ID
	ppprc.Position = column.Position
	ppprc.Cols = column.Cols
	if column.Xxl.Valid {
		ppprc.Xxl = &column.Xxl.Int16
	}
	if column.Xl.Valid {
		ppprc.Xl = &column.Xl.Int16
	}
	if column.Lg.Valid {
		ppprc.Lg = &column.Lg.Int16
	}
	if column.Md.Valid {
		ppprc.Md = &column.Md.Int16
	}
	if column.Sm.Valid {
		ppprc.Sm = &column.Sm.Int16
	}
	if column.Xs.Valid {
		ppprc.Xs = &column.Xs.Int16
	}
	if column.Offset.Valid {
		ppprc.Offset = &column.Offset.Int16
	}
	if column.OffsetXxl.Valid {
		ppprc.OffsetXxl = &column.OffsetXxl.Int16
	}
	if column.OffsetXl.Valid {
		ppprc.OffsetXl = &column.OffsetXl.Int16
	}
	if column.OffsetLg.Valid {
		ppprc.OffsetLg = &column.OffsetLg.Int16
	}
	if column.OffsetMd.Valid {
		ppprc.OffsetMd = &column.OffsetMd.Int16
	}
	if column.OffsetSm.Valid {
		ppprc.OffsetSm = &column.OffsetSm.Int16
	}
	if column.Order.Valid {
		ppprc.Order = &column.Order.Int16
	}
	if column.OrderXxl.Valid {
		ppprc.OrderXxl = &column.OrderXxl.Int16
	}
	if column.OrderXl.Valid {
		ppprc.OrderXl = &column.OrderXl.Int16
	}
	if column.OrderLg.Valid {
		ppprc.OrderLg = &column.OrderLg.Int16
	}
	if column.OrderMd.Valid {
		ppprc.OrderMd = &column.OrderMd.Int16
	}
	if column.OrderSm.Valid {
		ppprc.OrderSm = &column.OrderSm.Int16
	}
	if column.AlignSelf.Valid {
		ppprc.AlignSelf = &column.AlignSelf.String
	}
	if column.Content.Valid {
		ppprc.Content = &column.Content.String
	}
	if column.Module != nil {
		var module PublishedModule
		module.SetModule(column.Module)
		ppprc.Module = &module
	}
}
