package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type PagePartialRowColumn struct {
	gorm.Model
	PagePartialRowID uint `gorm:"not null"`
	ModuleID         sql.Null[uint]
	Cols             string `gorm:"no null;size:32"`
	Xxl              sql.NullInt16
	Xl               sql.NullInt16
	Lg               sql.NullInt16
	Md               sql.NullInt16
	Sm               sql.NullInt16
	Xs               sql.NullInt16
	Offset           sql.NullInt16
	OffsetXxl        sql.NullInt16
	OffsetXl         sql.NullInt16
	OffsetLg         sql.NullInt16
	OffsetMd         sql.NullInt16
	OffsetSm         sql.NullInt16
	Order            sql.NullInt16
	OrderXxl         sql.NullInt16
	OrderXl          sql.NullInt16
	OrderLg          sql.NullInt16
	OrderMd          sql.NullInt16
	OrderSm          sql.NullInt16
	AlignSelf        sql.NullString `gorm:"size:32"`
	Content          sql.NullString

	// Relationships.
	PagePartialRow PagePartialRow `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PagePartialRowID;references:ID"`
	Module         Module         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ModuleID;references:ID"`
}
