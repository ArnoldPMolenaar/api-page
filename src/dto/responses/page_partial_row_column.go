package responses

import (
	"api-page/main/src/models"
	"time"
)

type PagePartialRowColumn struct {
	ID        uint      `json:"id" `
	RowID     uint      `json:"rowId"`
	ModuleID  *uint     `json:"moduleId"`
	Position  uint      `json:"position"`
	Cols      string    `json:"cols"`
	Xxl       *int16    `json:"xxl"`
	Xl        *int16    `json:"xl"`
	Lg        *int16    `json:"lg"`
	Md        *int16    `json:"md"`
	Sm        *int16    `json:"sm"`
	Xs        *int16    `json:"xs"`
	Offset    *int16    `json:"offset"`
	OffsetXxl *int16    `json:"offsetXxl"`
	OffsetXl  *int16    `json:"offsetXl"`
	OffsetLg  *int16    `json:"offsetLg"`
	OffsetMd  *int16    `json:"offsetMd"`
	OffsetSm  *int16    `json:"offsetSm"`
	Order     *int16    `json:"order"`
	OrderXxl  *int16    `json:"orderXxl"`
	OrderXl   *int16    `json:"orderXl"`
	OrderLg   *int16    `json:"orderLg"`
	OrderMd   *int16    `json:"orderMd"`
	OrderSm   *int16    `json:"orderSm"`
	AlignSelf *string   `json:"alignSelf"`
	Content   *string   `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SetPagePartialRowColumn sets the PagePartialRowColumn response from the models.PagePartialRowColumn model.
func (pprc *PagePartialRowColumn) SetPagePartialRowColumn(column *models.PagePartialRowColumn) {
	pprc.ID = column.ID
	pprc.RowID = column.PagePartialRowID
	pprc.Position = column.Position
	pprc.Cols = column.Cols
	pprc.CreatedAt = column.CreatedAt
	pprc.UpdatedAt = column.UpdatedAt
	if column.ModuleID.Valid {
		pprc.ModuleID = &column.ModuleID.V
	}
	if column.Xxl.Valid {
		pprc.Xxl = &column.Xxl.Int16
	}
	if column.Xl.Valid {
		pprc.Xl = &column.Xl.Int16
	}
	if column.Lg.Valid {
		pprc.Lg = &column.Lg.Int16
	}
	if column.Md.Valid {
		pprc.Md = &column.Md.Int16
	}
	if column.Sm.Valid {
		pprc.Sm = &column.Sm.Int16
	}
	if column.Xs.Valid {
		pprc.Xs = &column.Xs.Int16
	}
	if column.Offset.Valid {
		pprc.Offset = &column.Offset.Int16
	}
	if column.OffsetXxl.Valid {
		pprc.OffsetXxl = &column.OffsetXxl.Int16
	}
	if column.OffsetXl.Valid {
		pprc.OffsetXl = &column.OffsetXl.Int16
	}
	if column.OffsetLg.Valid {
		pprc.OffsetLg = &column.OffsetLg.Int16
	}
	if column.OffsetMd.Valid {
		pprc.OffsetMd = &column.OffsetMd.Int16
	}
	if column.OffsetSm.Valid {
		pprc.OffsetSm = &column.OffsetSm.Int16
	}
	if column.Order.Valid {
		pprc.Order = &column.Order.Int16
	}
	if column.OrderXxl.Valid {
		pprc.OrderXxl = &column.OrderXxl.Int16
	}
	if column.OrderXl.Valid {
		pprc.OrderXl = &column.OrderXl.Int16
	}
	if column.OrderLg.Valid {
		pprc.OrderLg = &column.OrderLg.Int16
	}
	if column.OrderMd.Valid {
		pprc.OrderMd = &column.OrderMd.Int16
	}
	if column.OrderSm.Valid {
		pprc.OrderSm = &column.OrderSm.Int16
	}
	if column.AlignSelf.Valid {
		pprc.AlignSelf = &column.AlignSelf.String
	}
	if column.Content.Valid {
		pprc.Content = &column.Content.String
	}
}
