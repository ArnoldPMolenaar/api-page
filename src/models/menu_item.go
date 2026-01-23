package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type MenuItem struct {
	gorm.Model
	EnabledAt sql.NullTime
	VersionID uint   `gorm:"not null"`
	Name      string `gorm:"not null"`
	Icon      sql.NullString

	// Relationships.
	Version  Version            `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:VersionID;references:ID"`
	Indexing []MenuItemIndexing `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MenuItemID;references:ID"`
}
