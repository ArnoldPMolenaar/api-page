package responses

import (
	"api-page/main/src/models"
	"time"
)

type FooterRowColumn struct {
	ID        uint        `json:"id" `
	RowID     uint        `json:"rowId"`
	ModuleID  *uint       `json:"moduleId"`
	Position  uint        `json:"position"`
	Cols      string      `json:"cols"`
	Xxl       *int16      `json:"xxl"`
	Xl        *int16      `json:"xl"`
	Lg        *int16      `json:"lg"`
	Md        *int16      `json:"md"`
	Sm        *int16      `json:"sm"`
	Xs        *int16      `json:"xs"`
	Offset    *int16      `json:"offset"`
	OffsetXxl *int16      `json:"offsetXxl"`
	OffsetXl  *int16      `json:"offsetXl"`
	OffsetLg  *int16      `json:"offsetLg"`
	OffsetMd  *int16      `json:"offsetMd"`
	OffsetSm  *int16      `json:"offsetSm"`
	Order     *int16      `json:"order"`
	OrderXxl  *int16      `json:"orderXxl"`
	OrderXl   *int16      `json:"orderXl"`
	OrderLg   *int16      `json:"orderLg"`
	OrderMd   *int16      `json:"orderMd"`
	OrderSm   *int16      `json:"orderSm"`
	AlignSelf *string     `json:"alignSelf"`
	Content   *string     `json:"content"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Rows      []FooterRow `json:"rows"`
}

// SetFooterRowColumn sets the FooterRowColumn response from the models.FooterRowColumn model.
func (frc *FooterRowColumn) SetFooterRowColumn(column *models.FooterRowColumn) {
	frc.ID = column.ID
	frc.RowID = column.FooterRowID
	frc.Position = column.Position
	frc.Cols = column.Cols
	frc.CreatedAt = column.CreatedAt
	frc.UpdatedAt = column.UpdatedAt
	if column.ModuleID.Valid {
		frc.ModuleID = &column.ModuleID.V
	}
	if column.Xxl.Valid {
		frc.Xxl = &column.Xxl.Int16
	}
	if column.Xl.Valid {
		frc.Xl = &column.Xl.Int16
	}
	if column.Lg.Valid {
		frc.Lg = &column.Lg.Int16
	}
	if column.Md.Valid {
		frc.Md = &column.Md.Int16
	}
	if column.Sm.Valid {
		frc.Sm = &column.Sm.Int16
	}
	if column.Xs.Valid {
		frc.Xs = &column.Xs.Int16
	}
	if column.Offset.Valid {
		frc.Offset = &column.Offset.Int16
	}
	if column.OffsetXxl.Valid {
		frc.OffsetXxl = &column.OffsetXxl.Int16
	}
	if column.OffsetXl.Valid {
		frc.OffsetXl = &column.OffsetXl.Int16
	}
	if column.OffsetLg.Valid {
		frc.OffsetLg = &column.OffsetLg.Int16
	}
	if column.OffsetMd.Valid {
		frc.OffsetMd = &column.OffsetMd.Int16
	}
	if column.OffsetSm.Valid {
		frc.OffsetSm = &column.OffsetSm.Int16
	}
	if column.Order.Valid {
		frc.Order = &column.Order.Int16
	}
	if column.OrderXxl.Valid {
		frc.OrderXxl = &column.OrderXxl.Int16
	}
	if column.OrderXl.Valid {
		frc.OrderXl = &column.OrderXl.Int16
	}
	if column.OrderLg.Valid {
		frc.OrderLg = &column.OrderLg.Int16
	}
	if column.OrderMd.Valid {
		frc.OrderMd = &column.OrderMd.Int16
	}
	if column.OrderSm.Valid {
		frc.OrderSm = &column.OrderSm.Int16
	}
	if column.AlignSelf.Valid {
		frc.AlignSelf = &column.AlignSelf.String
	}
	if column.Content.Valid {
		frc.Content = &column.Content.String
	}

	frc.Rows = make([]FooterRow, len(column.FooterRows))
	for i := range column.FooterRows {
		row := FooterRow{}
		row.SetFooterRow(&column.FooterRows[i])
		frc.Rows[i] = row
	}
}
