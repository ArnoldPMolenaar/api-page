package models

import (
	"api-page/main/src/enums"
	"database/sql"
)

type PageIndexing struct {
	MenuItemID uint           `gorm:"primaryKey:true;autoIncrement:false;index:idx_page_indexing_option,unique"`
	Locale     string         `gorm:"primaryKey:true;autoIncrement:false;size:32;index:idx_page_indexing_option,unique"`
	Option     enums.Indexing `gorm:"not null;type:indexing;default:index;index:idx_page_indexing_option,unique"`
	Value      sql.NullString

	// Relationships.
	MenuItem MenuItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MenuItemID;references:ID"`
}
