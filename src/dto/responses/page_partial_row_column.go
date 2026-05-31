package responses

import (
	"api-page/main/src/models"
	"time"

	"github.com/ArnoldPMolenaar/api-utils/utils"
)

type PagePartialRowColumn struct {
	ID        uint             `json:"id" `
	RowID     uint             `json:"rowId"`
	ModuleID  *uint            `json:"moduleId"`
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
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
	Rows      []PagePartialRow `json:"rows"`
}

// SetPagePartialRowColumn sets the PagePartialRowColumn response from the models.PagePartialRowColumn model.
func (pprc *PagePartialRowColumn) SetPagePartialRowColumn(column *models.PagePartialRowColumn) {
	pprc.ID = column.ID
	pprc.RowID = column.PagePartialRowID
	pprc.Position = column.Position
	pprc.Cols = column.Cols
	pprc.CreatedAt = column.CreatedAt
	pprc.UpdatedAt = column.UpdatedAt
	pprc.ModuleID = utils.PtrFromNull[uint](column.ModuleID)
	pprc.Xxl = utils.PtrFromNullInt16(column.Xxl)
	pprc.Xl = utils.PtrFromNullInt16(column.Xl)
	pprc.Lg = utils.PtrFromNullInt16(column.Lg)
	pprc.Md = utils.PtrFromNullInt16(column.Md)
	pprc.Sm = utils.PtrFromNullInt16(column.Sm)
	pprc.Xs = utils.PtrFromNullInt16(column.Xs)
	pprc.Offset = utils.PtrFromNullInt16(column.Offset)
	pprc.OffsetXxl = utils.PtrFromNullInt16(column.OffsetXxl)
	pprc.OffsetXl = utils.PtrFromNullInt16(column.OffsetXl)
	pprc.OffsetLg = utils.PtrFromNullInt16(column.OffsetLg)
	pprc.OffsetMd = utils.PtrFromNullInt16(column.OffsetMd)
	pprc.OffsetSm = utils.PtrFromNullInt16(column.OffsetSm)
	pprc.Order = utils.PtrFromNullInt16(column.Order)
	pprc.OrderXxl = utils.PtrFromNullInt16(column.OrderXxl)
	pprc.OrderXl = utils.PtrFromNullInt16(column.OrderXl)
	pprc.OrderLg = utils.PtrFromNullInt16(column.OrderLg)
	pprc.OrderMd = utils.PtrFromNullInt16(column.OrderMd)
	pprc.OrderSm = utils.PtrFromNullInt16(column.OrderSm)
	pprc.AlignSelf = utils.PtrFromNullString(column.AlignSelf)
	pprc.Content = utils.PtrFromNullString(column.Content)

	pprc.Rows = make([]PagePartialRow, len(column.PagePartialRows))
	for i := range column.PagePartialRows {
		row := PagePartialRow{}
		row.SetPagePartialRow(&column.PagePartialRows[i])
		pprc.Rows[i] = row
	}
}
