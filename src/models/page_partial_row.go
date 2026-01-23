package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type PagePartialRow struct {
	gorm.Model
	PagePartialID   uint           `gorm:"not null"`
	NoGutters       bool           `gorm:"not null;default:false"`
	Dense           bool           `gorm:"not null;default:false"`
	Align           sql.NullString `gorm:"size:32"`
	AlignXxl        sql.NullString `gorm:"size:32"`
	AlignXl         sql.NullString `gorm:"size:32"`
	AlignLg         sql.NullString `gorm:"size:32"`
	AlignMd         sql.NullString `gorm:"size:32"`
	AlignSm         sql.NullString `gorm:"size:32"`
	AlignContent    sql.NullString `gorm:"size:32"`
	AlignContentXxl sql.NullString `gorm:"size:32"`
	AlignContentXl  sql.NullString `gorm:"size:32"`
	AlignContentLg  sql.NullString `gorm:"size:32"`
	AlignContentMd  sql.NullString `gorm:"size:32"`
	AlignContentSm  sql.NullString `gorm:"size:32"`
	Justify         sql.NullString `gorm:"size:32"`
	JustifyXxl      sql.NullString `gorm:"size:32"`
	JustifyXl       sql.NullString `gorm:"size:32"`
	JustifyLg       sql.NullString `gorm:"size:32"`
	JustifyMd       sql.NullString `gorm:"size:32"`
	JustifySm       sql.NullString `gorm:"size:32"`

	// Relationships.
	PagePartial PagePartial            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PagePartialID;references:ID"`
	Columns     []PagePartialRowColumn `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PagePartialRowID;references:ID"`
}
