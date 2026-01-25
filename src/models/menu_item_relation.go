package models

import "database/sql"

type MenuItemRelation struct {
	MenuID           uint           `gorm:"primaryKey:true;autoIncrement:false;index:idx_menu_item,unique"`
	MenuItemParentID sql.Null[uint] `gorm:"index:idx_menu_item,unique"`
	MenuItemChildID  uint           `gorm:"primaryKey:true;autoIncrement:false;index:idx_menu_item,unique"`
	Position         uint           `gorm:"not null;default:0;index:idx_menu_item_position,sort:asc;index:idx_menu_item,unique"`

	// Relationships.
	Menu           Menu     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MenuID;references:ID"`
	MenuItemParent MenuItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MenuItemParentID;references:ID"`
	MenuItemChild  MenuItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:MenuItemChildID;references:ID"`
}
