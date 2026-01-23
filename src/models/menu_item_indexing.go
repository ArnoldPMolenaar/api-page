package models

import (
	"api-page/main/src/enums"
	"database/sql"
)

type MenuItemIndexing struct {
	MenuItemID uint           `gorm:"primaryKey:true;autoIncrement:false"`
	Option     enums.Indexing `gorm:"primaryKey:true;type:indexing;default:index"`
	Value      sql.NullString

	// Relationships.
	MenuItem MenuItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MenuItemID;references:ID"`
}
