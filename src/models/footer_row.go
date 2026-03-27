package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type FooterRow struct {
	gorm.Model
	VersionID       uint   `gorm:"not null"`
	Locale          string `gorm:"not null;size:32"`
	Position        uint   `gorm:"not null"`
	NoGutters       bool   `gorm:"not null;default:false"`
	Dense           bool   `gorm:"not null;default:false"`
	Hashtag         sql.NullString
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
	Version Version           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:VersionID;references:ID"`
	Columns []FooterRowColumn `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:FooterRowID;references:ID"`
}
