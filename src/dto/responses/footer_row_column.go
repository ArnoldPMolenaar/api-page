package responses

import (
	"api-page/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
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
	frc.ModuleID = utils.PtrFromNull[uint](column.ModuleID)
	frc.Xxl = utils.PtrFromNullInt16(column.Xxl)
	frc.Xl = utils.PtrFromNullInt16(column.Xl)
	frc.Lg = utils.PtrFromNullInt16(column.Lg)
	frc.Md = utils.PtrFromNullInt16(column.Md)
	frc.Sm = utils.PtrFromNullInt16(column.Sm)
	frc.Xs = utils.PtrFromNullInt16(column.Xs)
	frc.Offset = utils.PtrFromNullInt16(column.Offset)
	frc.OffsetXxl = utils.PtrFromNullInt16(column.OffsetXxl)
	frc.OffsetXl = utils.PtrFromNullInt16(column.OffsetXl)
	frc.OffsetLg = utils.PtrFromNullInt16(column.OffsetLg)
	frc.OffsetMd = utils.PtrFromNullInt16(column.OffsetMd)
	frc.OffsetSm = utils.PtrFromNullInt16(column.OffsetSm)
	frc.Order = utils.PtrFromNullInt16(column.Order)
	frc.OrderXxl = utils.PtrFromNullInt16(column.OrderXxl)
	frc.OrderXl = utils.PtrFromNullInt16(column.OrderXl)
	frc.OrderLg = utils.PtrFromNullInt16(column.OrderLg)
	frc.OrderMd = utils.PtrFromNullInt16(column.OrderMd)
	frc.OrderSm = utils.PtrFromNullInt16(column.OrderSm)
	frc.AlignSelf = utils.PtrFromNullString(column.AlignSelf)
	frc.Content = utils.PtrFromNullString(column.Content)

	frc.Rows = make([]FooterRow, len(column.FooterRows))
	for i := range column.FooterRows {
		row := FooterRow{}
		row.SetFooterRow(&column.FooterRows[i])
		frc.Rows[i] = row
	}
}
